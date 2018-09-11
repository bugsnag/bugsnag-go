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
	"github.com/bugsnag/bugsnag-go/sessions"
)

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
// report payloads published to ts.URL will be put on the returned channel
func setup() (*httptest.Server, chan []byte) {
	reports := make(chan []byte, 10)
	sessionTracker = &testSessionTracker{}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		reports <- body
	})), reports
}

type testSessionTracker struct{}

func (t *testSessionTracker) StartSession(context.Context) context.Context {
	return context.Background()
}

func (t *testSessionTracker) GetSession(context.Context) *sessions.Session {
	return &sessions.Session{}
}

func TestConfigure(t *testing.T) {
	Configure(Configuration{
		APIKey: testAPIKey,
	})

	if Config.APIKey != testAPIKey {
		t.Errorf("Setting APIKey didn't work")
	}

	if New().Config.APIKey != testAPIKey {
		t.Errorf("Setting APIKey didn't work for new notifiers")
	}
}

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

	Notify(
		StartSession(context.Background()),
		fmt.Errorf("hello world"),
		generateSampleConfig(ts.URL),
		User{Id: "123", Name: "Conrad", Email: "me@cirw.in"},
		Context{"testing"},
		MetaData{"test": {
			"password": "sneaky",
			"value":    "able",
			"broken":   complex(1, 2),
			"recurse":  recurse,
		}},
	)

	json, err := simplejson.NewJson(<-reports)

	if err != nil {
		t.Fatal(err)
	}

	event := json.Get("events").GetIndex(0)

	assertPayload(t, json, eventJSON{
		App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
		Context:        "testing",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "lol",
		Session:        &sessionJSON{Events: eventCountsJSON{Handled: 0, Unhandled: 1}},
		Severity:       "warning",
		SeverityReason: &severityReasonJSON{Attributes: &severityAttributesJSON{Framework: ""}, Type: SeverityReasonHandledError},
		Unhandled:      false,
		User:           &User{Id: "123", Name: "Conrad", Email: "me@cirw.in"},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "hello world"}},
	})
	assertValidSession(t, event, false)

	for k, exp := range map[string]string{
		"metaData.test.password":        "[REDACTED]",
		"metaData.test.value":           "able",
		"metaData.test.broken":          "[complex128]",
		"metaData.test.recurse.Recurse": "[RECURSION]",
	} {
		if got := getString(event, k); got != exp {
			t.Errorf("Expected %s to be '%s' but was '%s'", k, exp, got)
		}
	}

	exception := event.Get("exceptions").GetIndex(0)
	checkFrame(t, getIndex(exception, "stacktrace", 0), stackFrame{File: "bugsnag_test.go", Method: "TestNotify", LineNumber: 93, InProject: true})
	checkFrame(t, getIndex(exception, "stacktrace", 1), stackFrame{File: "testing/testing.go", Method: "tRunner", InProject: false})
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

	go (&http.Server{
		Addr:     l.Addr().String(),
		Handler:  Handler(mux, generateSampleConfig(ts.URL), SeverityInfo),
		ErrorLog: log.New(ioutil.Discard, log.Prefix(), 0),
	}).Serve(l)

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
		Session:        &sessionJSON{Events: eventCountsJSON{Handled: 0, Unhandled: 1}},
		Severity:       "info",
		SeverityReason: &severityReasonJSON{Attributes: &severityAttributesJSON{Framework: ""}, Type: SeverityReasonHandledPanic},
		Unhandled:      true,
		User:           &User{Id: "127.0.0.1", Name: "", Email: ""},
		Exceptions:     []exceptionJSON{{ErrorClass: "runtime.plainError", Message: "send on closed channel"}},
	})
	event := getIndex(json, "events", 0)
	assertValidSession(t, event, true)
	for k, exp := range map[string]string{
		"metaData.request.httpMethod": "GET",
		"metaData.request.url":        "http://" + l.Addr().String() + "/ok?foo=bar",
	} {
		if got := getString(event, k); got != exp {
			t.Errorf("Expected %s to be '%s' but was '%s'", k, exp, got)
		}
	}
	for k, exp := range map[string]string{
		"metaData.request.params.foo":              "bar",
		"metaData.request.headers.Accept-Encoding": "gzip",
	} {
		if got := getFirstString(event, k); got != exp {
			t.Errorf("Expected %s to be '%s' but was '%s'", k, exp, got)
		}
	}

	exception := getIndex(event, "exceptions", 0)
	checkFrame(t, getIndex(exception, "stacktrace", 0), stackFrame{File: "runtime/panic.go", Method: "gopanic", InProject: false})
	checkFrame(t, getIndex(exception, "stacktrace", 3), stackFrame{File: "bugsnag_test.go", Method: "crashyHandler", LineNumber: 24, InProject: true})
}

func TestAutoNotify(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	var panicked error

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
		defer AutoNotify(StartSession(context.Background()), generateSampleConfig(ts.URL))

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
		Session:        &sessionJSON{Events: eventCountsJSON{Handled: 0, Unhandled: 1}},
		Severity:       "error",
		SeverityReason: &severityReasonJSON{Attributes: &severityAttributesJSON{Framework: ""}, Type: SeverityReasonHandledPanic}, //TODO: this should be unhandled panic!
		Unhandled:      true,
		User:           &User{},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "eggs"}},
	})
}

func TestRecover(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	var panicked interface{}

	func() {
		defer func() {
			panicked = recover()
		}()
		defer Recover(StartSession(context.Background()), generateSampleConfig(ts.URL))

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
		Session:        &sessionJSON{Events: eventCountsJSON{Handled: 0, Unhandled: 1}},
		Severity:       "warning",
		SeverityReason: &severityReasonJSON{Attributes: &severityAttributesJSON{Framework: ""}, Type: SeverityReasonHandledPanic},
		Unhandled:      false,
		User:           &User{},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "ham"}},
	})
}

//TODO: What does this test add over TestNotify?
func TestSeverityReasonNotifyErr(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	Notify(StartSession(context.Background()), fmt.Errorf("hello world"), generateSampleConfig(ts.URL))

	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}

	assertPayload(t, json, eventJSON{
		App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
		Context:        "",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "",
		Session:        &sessionJSON{Events: eventCountsJSON{Handled: 0, Unhandled: 1}},
		Severity:       "warning",
		SeverityReason: &severityReasonJSON{Attributes: &severityAttributesJSON{Framework: ""}, Type: SeverityReasonHandledError},
		Unhandled:      false,
		User:           &User{},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "hello world"}},
	})
}

func TestSeverityReasonNotifyCallback(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	OnBeforeNotify(func(event *Event, config *Configuration) error {
		event.Severity = SeverityInfo
		return nil
	})

	Notify(StartSession(context.Background()), fmt.Errorf("hello world"), generateSampleConfig(ts.URL))

	json, _ := simplejson.NewJson(<-reports)
	assertPayload(t, json, eventJSON{
		App:            &appJSON{ReleaseStage: "test", Type: "foo", Version: "1.2.3"},
		Context:        "",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "",
		Session:        &sessionJSON{Events: eventCountsJSON{Handled: 0, Unhandled: 1}},
		Severity:       "info",
		SeverityReason: &severityReasonJSON{Attributes: &severityAttributesJSON{Framework: ""}, Type: SeverityReasonCallbackSpecified},
		Unhandled:      false,
		User:           &User{},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "hello world"}},
	})
}

func TestConfigureTwice(t *testing.T) {
	sessionTracker = nil

	l := logger{}
	Configure(Configuration{Logger: &l})
	if l.msg != "" {
		t.Errorf("unexpected log message: %s", l.msg)
	}
	Configure(Configuration{})
	if got, exp := l.msg, configuredMultipleTimes; exp != got {
		t.Errorf("unexpected log message: '%s', expected '%s'", got, exp)
	}
}

func generateSampleConfig(endpoint string) Configuration {
	return Configuration{
		APIKey:          testAPIKey,
		Endpoints:       Endpoints{Notify: endpoint},
		ProjectPackages: []string{"github.com/bugsnag/bugsnag-go"},
		Logger:          log.New(ioutil.Discard, log.Prefix(), log.Flags()),
		ReleaseStage:    "test",
		AppType:         "foo",
		AppVersion:      "1.2.3",
		Hostname:        "web1",
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

	event := report.Get("events").GetIndex(0)
	exception := getIndex(event, "exceptions", 0)

	for _, tc := range []struct {
		prop     string
		exp, got interface{}
	}{
		{prop: "API Key", exp: testAPIKey, got: getString(report, "apiKey")},

		{prop: "notifier name", exp: "Bugsnag Go", got: getString(report, "notifier.name")},
		{prop: "notifier version", exp: VERSION, got: getString(report, "notifier.version")},
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
		{prop: "payload version", exp: "2", got: getString(event, "payloadVersion")},

		{prop: "severity", exp: exp.Severity, got: getString(event, "severity")},

		{
			prop: "severity reason attribute: 'framework'",
			exp:  exp.SeverityReason.Attributes.Framework,
			got:  getString(event, "severityReason.attributes.framework"),
		},

		{
			prop: "severity reason type",
			exp:  string(exp.SeverityReason.Type),
			got:  getString(event, "severityReason.type"),
		},
	} {
		if tc.got != tc.exp {
			t.Errorf("Wrong %s: expected '%v' but got '%v'", tc.prop, tc.exp, tc.got)
		}
	}
}

func assertValidSession(t *testing.T, event *simplejson.Json, unhandled bool) {
	if sessionID := event.GetPath("session", "id").MustString(); len(sessionID) != 36 {
		t.Errorf("Expected a valid session ID to be set but was '%s'", sessionID)
	}
	if _, e := time.Parse(time.RFC3339, event.GetPath("session", "startedAt").MustString()); e != nil {
		t.Error(e)
	}
	expHandled, expUnhandled := 1, 0
	if unhandled {
		expHandled, expUnhandled = expUnhandled, expHandled
	}
	if got := event.GetPath("session", "events", "unhandled").MustInt(); got != expUnhandled {
		t.Errorf("Expected %d unhandled events in session but was %d", expUnhandled, got)
	}
	if got := event.GetPath("session", "events", "handled").MustInt(); got != expHandled {
		t.Errorf("Expected %d handled events in session but was %d", expHandled, got)
	}
}

func checkFrame(t *testing.T, frame *simplejson.Json, exp stackFrame) {
	if got := getString(frame, "file"); got != exp.File {
		t.Errorf("Expected frame file to be '%s' but was '%s'", exp.File, got)
	}
	if got := getString(frame, "method"); got != exp.Method {
		t.Errorf("Expected frame method to be '%s' but was '%s'", exp.Method, got)
	}
	if got := getInt(frame, "lineNumber"); got != exp.LineNumber && exp.InProject { // Don't check files that vary per version of go
		t.Errorf("Expected frame line number to be %d but was %d", exp.LineNumber, got)
	}
	if got := getBool(frame, "inProject"); got != exp.InProject {
		t.Errorf("Expected frame inProject to be '%v' but was '%v'", exp.InProject, got)
	}
}
