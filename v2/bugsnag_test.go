package bugsnag

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go/v2/sessions"
)

// ATTENTION - tests in this file are changing global state variables
// like default config or default report publisher
// TAKE CARE to reset them to default after testcase!

// The line numbers of this method are used in tests.
// If you move this function you'll have to change tests
func crashyHandler(w http.ResponseWriter, r *http.Request) {
	c := make(chan int)
	close(c)
	c <- 1
}

type _recurse struct {
	Recurse *_recurse
}

const (
	unhandled = true
	handled   = false
)

var testAPIKey = "166f5ad3590596f9aa8d601ea89af845"

type logger struct{ msg string }

func (l *logger) Printf(format string, v ...interface{}) { l.msg = format }

// setup sets up a simple sessionTracker and returns a test event server for receiving the event payloads.
// report payloads published to the returned server's URL will be put on the returned channel
func setup() (*httptest.Server, chan []byte) {
	reports := make(chan []byte, 10)
	sessionTracker = &testSessionTracker{}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		reports <- body
	})), reports
}

// setupState resets global state so that we can run tests independently
// Using this is replacing the call to "Configure"
// Context needed to stop the publisher's delivery goroutine
func setupState() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	// Independent of the default one
	publisher = newPublisher()
	go publisher.delivery()
	publisher.setMainProgramContext(ctx)

	return ctx, cancel
}

type testSessionTracker struct{}

func (t *testSessionTracker) StartSession(context.Context) context.Context {
	return context.Background()
}

func (t *testSessionTracker) IncrementEventCountAndGetSession(context.Context, bool) *sessions.Session {
	return &sessions.Session{}
}

func (t *testSessionTracker) FlushSessions() {}

func TestNotify(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()
	sessionTracker = nil
	startSessionTracking()

	recurse := _recurse{}
	recurse.Recurse = &recurse

	OnBeforeNotify(func(event *Event, config *Configuration) error {
		if event.Context == "testing" {
			event.GroupingHash = "lol"
		}
		return nil
	})

	md := MetaData{"test": {"password": "sneaky", "value": "able", "broken": complex(1, 2), "recurse": recurse}}
	user := User{Id: "123", Name: "Conrad", Email: "me@cirw.in"}

	ctx, cancel := setupState()
	defer cancel()
	config := generateSampleConfig(ts.URL, ctx)

	Notify(fmt.Errorf("hello world"), StartSession(context.Background()), config, user, ErrorClass{Name: "ExpectedErrorClass"}, Context{"testing"}, md)

	json, err := simplejson.NewJson(<-reports)

	if err != nil {
		t.Fatal(err)
	}

	event := getIndex(json, "events", 0)

	assertPayload(t, json, eventJSON{
		App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
		Context:        "testing",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "lol",
		Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 0, Unhandled: 1}},
		Severity:       "warning",
		SeverityReason: &severityReasonJSON{Type: SeverityReasonHandledError},
		Unhandled:      false,
		Request:        &RequestJSON{},
		User:           &User{Id: "123", Name: "Conrad", Email: "me@cirw.in"},
		Exceptions:     []exceptionJSON{{ErrorClass: "ExpectedErrorClass", Message: "hello world"}},
	})
	assertValidSession(t, event, handled)

	for k, exp := range map[string]string{
		"metaData.test.password":        "[FILTERED]",
		"metaData.test.value":           "able",
		"metaData.test.broken":          "[complex128]",
		"metaData.test.recurse.Recurse": "[RECURSION]",
	} {
		if got := getString(event, k); got != exp {
			t.Errorf("Expected %s to be '%s' but was '%s'", k, exp, got)
		}
	}

	exception := getIndex(event, "exceptions", 0)
	verifyExistsInStackTrace(t, exception, &StackFrame{File: "bugsnag_test.go", Method: "TestNotify", LineNumber: 105, InProject: true})
}

type testPublisher struct {
	sync bool
}

func (tp *testPublisher) publishReport(p *payload) error {
	tp.sync = p.Synchronous
	return nil
}

func (tp *testPublisher) setMainProgramContext(context.Context) {
}

func (tp *testPublisher) delivery() {
}

func TestNotifySyncThenAsync(t *testing.T) {
	ts, _ := setup()
	defer ts.Close()

	_, cancel := setupState()
	defer cancel()

	pub := new(testPublisher)
	publisher = pub
	defer func() { publisher = newPublisher() }()

	Notify(fmt.Errorf("oopsie"))
	if pub.sync {
		t.Errorf("Expected notify to be async by default")
	}

	defaultNotifier.NotifySync(fmt.Errorf("oopsie"), true)
	if !pub.sync {
		t.Errorf("Expected notify to be sent synchronously when calling NotifySync with true")
	}

	Notify(fmt.Errorf("oopsie"))
	if pub.sync {
		t.Errorf("Expected notify to be sent asynchronously when calling Notify regardless of previous NotifySync call")
	}
}

func TestHandlerFunc(t *testing.T) {
	eventserver, reports := setup()
	defer eventserver.Close()

	// Save old config to restore later after the tests
	oldConfig := Config.clone()
	config := generateSampleConfig(eventserver.URL, context.Background())
	Config.update(&config)

	// NOTE - this testcase will print a panic in verbose mode
	t.Run("unhandled", func(st *testing.T) {
		_, cancel := setupState()
		defer cancel()

		sessionTracker = nil
		startSessionTracking()
		ts := httptest.NewServer(HandlerFunc(crashyHandler))
		defer ts.Close()

		http.Get(ts.URL + "/unhandled")

		json, _ := simplejson.NewJson(<-reports)
		assertPayload(t, json, eventJSON{
			App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
			Context:        "/unhandled",
			Device:         &deviceJSON{Hostname: "web1"},
			GroupingHash:   "",
			Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 0, Unhandled: 1}},
			Severity:       "error",
			SeverityReason: &severityReasonJSON{Type: SeverityReasonHandledPanic},
			Unhandled:      true,
			Request: &RequestJSON{
				Headers:    map[string]string{"Accept-Encoding": "gzip"},
				HTTPMethod: "GET",
				URL:        ts.URL + "/unhandled",
			},
			User:       &User{Id: "127.0.0.1", Name: "", Email: ""},
			Exceptions: []exceptionJSON{{ErrorClass: "runtime.plainError", Message: "send on closed channel"}},
		})
		event := getIndex(json, "events", 0)
		if got, exp := getString(event, "request.headers.Accept-Encoding"), "gzip"; got != exp {
			st.Errorf("expected Accept-Encoding header to be '%s' but was '%s'", exp, got)
		}
		if got, exp := getString(event, "request.httpMethod"), "GET"; got != exp {
			st.Errorf("expected HTTP method to be '%s' but was '%s'", exp, got)
		}
		if got, exp := getString(event, "request.url"), "/unhandled"; !strings.Contains(got, exp) {
			st.Errorf("expected request URL to contain '%s' but was '%s'", exp, got)
		}
		assertValidSession(st, event, unhandled)
	})

	t.Run("handled", func(st *testing.T) {
		_, cancel := setupState()
		defer cancel()

		sessionTracker = nil
		startSessionTracking()
		ts := httptest.NewServer(HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Notify(fmt.Errorf("oopsie"), r.Context())
		}))
		defer ts.Close()

		http.Get(ts.URL + "/handled")

		json, _ := simplejson.NewJson(<-reports)
		assertPayload(t, json, eventJSON{
			App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
			Context:        "/handled",
			Device:         &deviceJSON{Hostname: "web1"},
			GroupingHash:   "",
			Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 1, Unhandled: 0}},
			Severity:       "warning",
			SeverityReason: &severityReasonJSON{Type: SeverityReasonHandledError},
			Unhandled:      false,
			Request: &RequestJSON{
				Headers:    map[string]string{"Accept-Encoding": "gzip"},
				HTTPMethod: "GET",
				URL:        ts.URL + "/handled",
			},
			User:       &User{Id: "127.0.0.1", Name: "", Email: ""},
			Exceptions: []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "oopsie"}},
		})
		event := getIndex(json, "events", 0)
		if got, exp := getString(event, "request.headers.Accept-Encoding"), "gzip"; got != exp {
			st.Errorf("expected Accept-Encoding header to be '%s' but was '%s'", exp, got)
		}
		if got, exp := getString(event, "request.httpMethod"), "GET"; got != exp {
			st.Errorf("expected HTTP method to be '%s' but was '%s'", exp, got)
		}
		if got, exp := getString(event, "request.url"), "/handled"; !strings.Contains(got, exp) {
			st.Errorf("expected request URL to contain '%s' but was '%s'", exp, got)
		}
		assertValidSession(st, event, handled)
	})

	Config = *oldConfig
}

func TestHandler(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/", crashyHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go (&http.Server{
		Addr:     l.Addr().String(),
		Handler:  Handler(mux, generateSampleConfig(ts.URL, ctx), SeverityInfo),
		ErrorLog: log.New(ioutil.Discard, log.Prefix(), 0),
	}).Serve(l)

	sessionTracker = nil
	startSessionTracking()

	http.Get("http://" + l.Addr().String() + "/ok?foo=bar")

	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}

	assertPayload(t, json, eventJSON{
		App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
		Context:        "/ok",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "",
		Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 0, Unhandled: 1}},
		Severity:       "info",
		SeverityReason: &severityReasonJSON{Type: SeverityReasonHandledPanic},
		Unhandled:      true,
		User:           &User{Id: "127.0.0.1", Name: "", Email: ""},
		Request: &RequestJSON{
			Headers:    map[string]string{"Accept-Encoding": "gzip"},
			HTTPMethod: "GET",
			URL:        "http://" + l.Addr().String() + "/ok?foo=bar",
		},
		Exceptions: []exceptionJSON{{ErrorClass: "runtime.plainError", Message: "send on closed channel"}},
	})
	event := getIndex(json, "events", 0)
	if got, exp := getString(event, "request.headers.Accept-Encoding"), "gzip"; got != exp {
		t.Errorf("expected Accept-Encoding header to be '%s' but was '%s'", exp, got)
	}
	if got, exp := getString(event, "request.httpMethod"), "GET"; got != exp {
		t.Errorf("expected HTTP method to be '%s' but was '%s'", exp, got)
	}
	if got, exp := getString(event, "request.url"), "/ok?foo=bar"; !strings.Contains(got, exp) {
		t.Errorf("expected request URL to be '%s' but was '%s'", exp, got)
	}
	assertValidSession(t, event, unhandled)
	if got, exp := getFirstString(event, "metaData.request.params.foo"), "bar"; got != exp {
		t.Errorf("Expected metadata params 'foo' to be '%s' but was '%s'", exp, got)
	}

	exception := getIndex(event, "exceptions", 0)
	verifyExistsInStackTrace(t, exception, &StackFrame{File: "bugsnag_test.go", Method: "crashyHandler", InProject: true, LineNumber: 28})
}

func TestAutoNotify(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	var panicked error
	ctx, cancel := setupState()
	defer cancel()

	func() {
		defer func() {
			p := recover()
			switch p.(type) {
			case error:
				panicked = p.(error)
			default:
				t.Fatalf("Unexpected panic happened. Expected 'eggs' Error but was a(n) <%T> with value <%+v>", p, p)
			}
		}()
		defer AutoNotify(StartSession(context.Background()), generateSampleConfig(ts.URL, ctx))

		panic(fmt.Errorf("eggs"))
	}()

	if panicked.Error() != "eggs" {
		t.Errorf("didn't re-panic")
	}

	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}

	assertPayload(t, json, eventJSON{
		App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
		Context:        "",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "",
		Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 0, Unhandled: 1}},
		Severity:       "error",
		SeverityReason: &severityReasonJSON{Type: SeverityReasonHandledPanic},
		Unhandled:      true,
		User:           &User{},
		Request:        &RequestJSON{},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "eggs"}},
	})
}

func TestRecover(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()
	ctx, cancel := setupState()
	defer cancel()

	var panicked interface{}

	func() {
		defer func() {
			panicked = recover()
		}()
		defer Recover(StartSession(context.Background()), generateSampleConfig(ts.URL, ctx))

		panic("ham")
	}()

	if panicked != nil {
		t.Errorf("Did not expect a panic but repanicked")
	}

	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}

	assertPayload(t, json, eventJSON{
		App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
		Context:        "",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "",
		Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 0, Unhandled: 1}},
		Severity:       "warning",
		SeverityReason: &severityReasonJSON{Type: SeverityReasonHandledPanic},
		Unhandled:      false,
		Request:        &RequestJSON{},
		User:           &User{},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "ham"}},
	})
}

func TestRecoverCustomHandledState(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()
	ctx, cancel := setupState()
	defer cancel()

	var panicked interface{}

	func() {
		defer func() {
			panicked = recover()
		}()
		handledState := HandledState{
			SeverityReason:   SeverityReasonHandledPanic,
			OriginalSeverity: SeverityError,
			Unhandled:        true,
		}
		defer Recover(handledState, StartSession(context.Background()), generateSampleConfig(ts.URL, ctx))

		panic("at the disco?")
	}()

	if panicked != nil {
		t.Errorf("Did not expect a panic but repanicked")
	}
	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}

	assertPayload(t, json, eventJSON{
		App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
		Context:        "",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "",
		Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 0, Unhandled: 1}},
		Severity:       "error",
		SeverityReason: &severityReasonJSON{Type: SeverityReasonHandledPanic},
		Unhandled:      true,
		Request:        &RequestJSON{},
		User:           &User{},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "at the disco?"}},
	})
}

func TestSeverityReasonNotifyCallback(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()
	ctx, cancel := setupState()
	defer cancel()

	OnBeforeNotify(func(event *Event, config *Configuration) error {
		event.Severity = SeverityInfo
		return nil
	})

	Notify(fmt.Errorf("hello world"), generateSampleConfig(ts.URL, ctx), StartSession(context.Background()))

	json, _ := simplejson.NewJson(<-reports)
	assertPayload(t, json, eventJSON{
		App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
		Context:        "",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "",
		Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 0, Unhandled: 1}},
		Severity:       "info",
		SeverityReason: &severityReasonJSON{Type: SeverityReasonCallbackSpecified},
		Unhandled:      false,
		Request:        &RequestJSON{},
		User:           &User{},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "hello world"}},
	})
}

func TestNotifyWithoutError(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()
	ctx, cancel := setupState()
	defer cancel()

	oldConfig := Config.clone()
	config := generateSampleConfig(ts.URL, ctx)
	config.Synchronous = true
	l := logger{}
	config.Logger = &l
	Config.update(&config)

	Notify(nil, StartSession(context.Background()), config)

	select {
	case r := <-reports:
		t.Fatalf("Unexpected request made to bugsnag: %+v", string(r))
	default:
		for _, exp := range []string{"ERROR", "error", "Bugsnag", "not notified"} {
			if got := l.msg; !strings.Contains(got, exp) {
				t.Errorf("Expected to see '%s' in logged message but logged message was '%s'", exp, got)
			}
		}
	}
	Config = *oldConfig
}

func TestConfigureTwice(t *testing.T) {
	publisher = newPublisher()
	Configure(Configuration{})
	if !Config.IsAutoCaptureSessions() {
		t.Errorf("Expected auto capture sessions to be enabled by default")
	}
	Configure(Configuration{AutoCaptureSessions: false})
	if Config.IsAutoCaptureSessions() {
		t.Errorf("Expected auto capture sessions to be disabled when configured")
	}
	Configure(Configuration{AutoCaptureSessions: true})
	if !Config.IsAutoCaptureSessions() {
		t.Errorf("Expected auto capture sessions to be enabled when configured")
	}
}

func generateSampleConfig(endpoint string, ctx context.Context) Configuration {
	return Configuration{
		APIKey:          testAPIKey,
		Endpoints:       Endpoints{Notify: endpoint},
		ProjectPackages: []string{"github.com/bugsnag/bugsnag-go"},
		Logger:          log.New(ioutil.Discard, log.Prefix(), log.Flags()),
		ReleaseStage:    "test",
		AppType:         "foo",
		AppVersion:      "1.2.3",
		Hostname:        "web1",
		MainContext:     ctx,
	}
}

func get(j *simplejson.Json, path string) *simplejson.Json {
	return j.GetPath(strings.Split(path, ".")...)
}
func getBool(j *simplejson.Json, path string) bool {
	return get(j, path).MustBool()
}
func getInt(j *simplejson.Json, path string) int {
	return get(j, path).MustInt()
}
func getString(j *simplejson.Json, path string) string {
	return get(j, path).MustString()
}
func getIndex(j *simplejson.Json, path string, index int) *simplejson.Json {
	return get(j, path).GetIndex(index)
}
func getFirstString(j *simplejson.Json, path string) string {
	return getIndex(j, path, 0).MustString()
}

// assertPayload compares the payload that was received by the event-server to
// the expected report JSON payload
func assertPayload(t *testing.T, report *simplejson.Json, exp eventJSON) {
	expException := exp.Exceptions[0]

	event := getIndex(report, "events", 0)
	exception := getIndex(event, "exceptions", 0)

	for _, tc := range []struct {
		prop     string
		exp, got interface{}
	}{
		{prop: "API Key", exp: testAPIKey, got: getString(report, "apiKey")},

		{prop: "notifier name", exp: "Bugsnag Go", got: getString(report, "notifier.name")},
		{prop: "notifier version", exp: Version, got: getString(report, "notifier.version")},
		{prop: "notifier url", exp: "https://github.com/bugsnag/bugsnag-go", got: getString(report, "notifier.url")},

		{prop: "exception message", exp: expException.Message, got: getString(exception, "message")},
		{prop: "exception error class", exp: expException.ErrorClass, got: getString(exception, "errorClass")},

		{prop: "unhandled", exp: exp.Unhandled, got: getBool(event, "unhandled")},

		{prop: "app version", exp: exp.App.Version, got: getString(event, "app.version")},
		{prop: "app release stage", exp: exp.App.ReleaseStage, got: getString(event, "app.releaseStage")},
		{prop: "app type", exp: exp.App.Type, got: getString(event, "app.type")},

		{prop: "user id", exp: exp.User.Id, got: getString(event, "user.id")},
		{prop: "user name", exp: exp.User.Name, got: getString(event, "user.name")},
		{prop: "user email", exp: exp.User.Email, got: getString(event, "user.email")},

		{prop: "context", exp: exp.Context, got: getString(event, "context")},
		{prop: "device hostname", exp: exp.Device.Hostname, got: getString(event, "device.hostname")},
		{prop: "grouping hash", exp: exp.GroupingHash, got: getString(event, "groupingHash")},
		{prop: "payload version", exp: "4", got: getString(event, "payloadVersion")},

		{prop: "severity", exp: exp.Severity, got: getString(event, "severity")},
		{prop: "severity reason type", exp: string(exp.SeverityReason.Type), got: getString(event, "severityReason.type")},

		{prop: "request header 'Accept-Encoding'", exp: string(exp.Request.Headers["Accept-Encoding"]), got: getString(event, "request.headers.Accept-Encoding")},
		{prop: "request HTTP method", exp: string(exp.Request.HTTPMethod), got: getString(event, "request.httpMethod")},
		{prop: "request URL", exp: string(exp.Request.URL), got: getString(event, "request.url")},
	} {
		if tc.got != tc.exp {
			t.Errorf("Wrong %s: expected '%v' but got '%v'", tc.prop, tc.exp, tc.got)
		}
	}
}

func assertValidSession(t *testing.T, event *simplejson.Json, unhandled bool) {
	if sessionID := getString(event, "session.id"); len(sessionID) != 36 {
		t.Errorf("Expected a valid session ID to be set but was '%s'", sessionID)
	}
	if _, e := time.Parse(time.RFC3339, getString(event, "session.startedAt")); e != nil {
		t.Error(e)
	}
	expHandled, expUnhandled := 1, 0
	if unhandled {
		expHandled, expUnhandled = expUnhandled, expHandled
	}
	if got := getInt(event, "session.events.unhandled"); got != expUnhandled {
		t.Errorf("Expected %d unhandled events in session but was %d", expUnhandled, got)
	}
	if got := getInt(event, "session.events.handled"); got != expHandled {
		t.Errorf("Expected %d handled events in session but was %d", expHandled, got)
	}
}

func verifyExistsInStackTrace(t *testing.T, exception *simplejson.Json, exp *StackFrame) {
	isFile := func(frame *simplejson.Json) bool { return strings.HasSuffix(getString(frame, "file"), exp.File) }
	isMethod := func(frame *simplejson.Json) bool { return getString(frame, "method") == exp.Method }
	isLineNumber := func(frame *simplejson.Json) bool { return getInt(frame, "lineNumber") == exp.LineNumber }

	arr, _ := exception.Get("stacktrace").Array()
	for i := 0; i < len(arr); i++ {
		frame := getIndex(exception, "stacktrace", i)
		if isFile(frame) && isMethod(frame) && isLineNumber(frame) {
			return
		}
	}
	t.Errorf("Could not find expected stackframe %v in exception '%v'", exp, exception)
}
