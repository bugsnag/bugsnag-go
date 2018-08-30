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

	"github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go/sessions"
)

type _recurse struct {
	Recurse *_recurse
}

var testAPIKey = "166f5ad3590596f9aa8d601ea89af845"

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

	recurse := _recurse{}
	recurse.Recurse = &recurse

	OnBeforeNotify(func(event *Event, config *Configuration) error {
		if event.Context == "testing" {
			event.GroupingHash = "lol"
		}
		return nil
	})

	Notify(fmt.Errorf("hello world"),
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

	if json.Get("apiKey").MustString() != testAPIKey {
		t.Errorf("Wrong api key in payload")
	}

	if json.GetPath("notifier", "name").MustString() != "Bugsnag Go" {
		t.Errorf("Wrong notifier name in payload")
	}

	event := json.Get("events").GetIndex(0)

	for k, value := range map[string]string{
		"payloadVersion":                "2",
		"severity":                      "warning",
		"context":                       "testing",
		"groupingHash":                  "lol",
		"app.releaseStage":              "test",
		"app.type":                      "foo",
		"app.version":                   "1.2.3",
		"device.hostname":               "web1",
		"user.id":                       "123",
		"user.name":                     "Conrad",
		"user.email":                    "me@cirw.in",
		"metaData.test.password":        "[REDACTED]",
		"metaData.test.value":           "able",
		"metaData.test.broken":          "[complex128]",
		"metaData.test.recurse.Recurse": "[RECURSION]",
	} {
		key := strings.Split(k, ".")
		if event.GetPath(key...).MustString() != value {
			t.Errorf("Wrong %v: %v != %v", key, event.GetPath(key...).MustString(), value)
		}
	}

	exception := event.Get("exceptions").GetIndex(0)

	if exception.Get("message").MustString() != "hello world" {
		t.Errorf("Wrong message in payload")
	}

	if exception.Get("errorClass").MustString() != "*errors.errorString" {
		t.Errorf("Wrong errorClass in payload: %v", exception.Get("errorClass").MustString())
	}

	frame0 := exception.Get("stacktrace").GetIndex(0)
	if frame0.Get("file").MustString() != "bugsnag_test.go" ||
		frame0.Get("method").MustString() != "TestNotify" ||
		frame0.Get("inProject").MustBool() != true ||
		frame0.Get("lineNumber").MustInt() == 0 {
		t.Errorf("Wrong frame0")
	}

	frame1 := exception.Get("stacktrace").GetIndex(1)

	if frame1.Get("file").MustString() != "testing/testing.go" ||
		frame1.Get("method").MustString() != "tRunner" ||
		frame1.Get("inProject").MustBool() != false ||
		frame1.Get("lineNumber").MustInt() == 0 {
		t.Errorf("Wrong frame1")
	}
}

func TestHandler(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", crashyHandler)
	srv := http.Server{
		Addr: l.Addr().String(),
		Handler: Handler(mux, Configuration{
			APIKey:          testAPIKey,
			Endpoints:       Endpoints{Notify: ts.URL},
			ProjectPackages: []string{"github.com/bugsnag/bugsnag-go"},
			Logger:          log.New(ioutil.Discard, log.Prefix(), log.Flags()),
		}, SeverityInfo),
		ErrorLog: log.New(ioutil.Discard, log.Prefix(), 0),
	}

	go srv.Serve(l)

	http.Get("http://" + l.Addr().String() + "/ok?foo=bar")
	l.Close()

	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}

	if json.Get("apiKey").MustString() != testAPIKey {
		t.Errorf("Wrong api key in payload")
	}

	if json.GetPath("notifier", "name").MustString() != "Bugsnag Go" {
		t.Errorf("Wrong notifier name in payload")
	}

	event := json.Get("events").GetIndex(0)

	for k, value := range map[string]string{
		"payloadVersion":              "2",
		"severity":                    "info",
		"user.id":                     "127.0.0.1",
		"metaData.request.url":        "http://" + l.Addr().String() + "/ok?foo=bar",
		"metaData.request.httpMethod": "GET",
	} {
		key := strings.Split(k, ".")
		if event.GetPath(key...).MustString() != value {
			t.Errorf("Wrong %v: %v != %v", key, event.GetPath(key...).MustString(), value)
		}
	}

	if event.GetPath("metaData", "request", "params", "foo").GetIndex(0).MustString() != "bar" {
		t.Errorf("missing GET params in request metadata")
	}

	if event.GetPath("metaData", "request", "headers", "Accept-Encoding").GetIndex(0).MustString() != "gzip" {
		t.Errorf("missing GET params in request metadata: %v", event.GetPath("metaData", "request", "headers"))
	}

	exception := event.Get("exceptions").GetIndex(0)

	if !strings.Contains(exception.Get("message").MustString(), "send on closed channel") {
		t.Errorf("Wrong message in payload: %v '%v'", exception.Get("message").MustString(), "runtime error: send on closed channel")
	}

	errorClass := exception.Get("errorClass").MustString()
	if errorClass != "runtime.errorCString" && errorClass != "*errors.errorString" && errorClass != "runtime.plainError" {
		t.Errorf("Wrong errorClass in payload: %v, expected '%v', '%v', '%v'",
			exception.Get("errorClass").MustString(),
			"runtime.errorCString", "*errors.errorString", "runtime.plainError")
	}

	frame0 := exception.Get("stacktrace").GetIndex(0)

	file0 := frame0.Get("file").MustString()
	if !strings.HasPrefix(file0, "runtime/panic") ||
		frame0.Get("inProject").MustBool() != false {
		t.Errorf("Wrong frame0: %v", frame0)
	}

	frame3 := exception.Get("stacktrace").GetIndex(3)

	if frame3.Get("file").MustString() != "bugsnag_test.go" ||
		frame3.Get("method").MustString() != "crashyHandler" ||
		frame3.Get("inProject").MustBool() != true ||
		frame3.Get("lineNumber").MustInt() == 0 {
		t.Errorf("Wrong frame3: %v", frame3)
	}
}

func TestAutoNotify(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	var panicked interface{}

	func() {
		defer func() {
			panicked = recover()
		}()
		defer AutoNotify(Configuration{Endpoints: Endpoints{Notify: ts.URL}, APIKey: testAPIKey})

		panic("eggs")
	}()

	// Note: If this line panics attempting to convert a `runtime.errorString`
	// into `string` then comment out the `panicked = recover()` line above, as
	// the panic you received here is not the one we expected.
	if panicked.(string) != "eggs" {
		t.Errorf("didn't re-panic")
	}

	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}

	event := json.Get("events").GetIndex(0)

	if event.Get("severity").MustString() != "error" {
		t.Errorf("severity should be error")
	}
	exception := event.Get("exceptions").GetIndex(0)

	if exception.Get("message").MustString() != "eggs" {
		t.Errorf("caught wrong panic")
	}
	assertSeverityReasonEqual(t, json, "error", "handledPanic", true)
}

func TestRecover(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	var panicked interface{}

	func() {
		defer func() {
			panicked = recover()
		}()
		defer Recover(Configuration{Endpoints: Endpoints{Notify: ts.URL}, APIKey: testAPIKey})

		panic("ham")
	}()

	if panicked != nil {
		t.Errorf("re-panick'd")
	}

	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}

	event := json.Get("events").GetIndex(0)

	if event.Get("severity").MustString() != "warning" {
		t.Errorf("severity should be warning")
	}
	exception := event.Get("exceptions").GetIndex(0)

	if exception.Get("message").MustString() != "ham" {
		t.Errorf("caught wrong panic")
	}
	assertSeverityReasonEqual(t, json, "warning", "handledPanic", false)
}

func TestSeverityReasonNotifyErr(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	Notify(fmt.Errorf("hello world"), generateSampleConfig(ts.URL))

	json, _ := simplejson.NewJson(<-reports)
	assertSeverityReasonEqual(t, json, "warning", "handledError", false)
}

func TestSeverityReasonNotifyCallback(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	OnBeforeNotify(func(event *Event, config *Configuration) error {
		event.Severity = SeverityInfo
		return nil
	})

	Notify(fmt.Errorf("hello world"), generateSampleConfig(ts.URL))

	json, _ := simplejson.NewJson(<-reports)
	assertSeverityReasonEqual(t, json, "info", "userCallbackSetSeverity", false)
}

type logger struct{ msg string }

func (l *logger) Printf(format string, v ...interface{}) { l.msg = format }

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

func assertSeverityReasonEqual(t *testing.T, json *simplejson.Json, expSeverity string, reasonType string, expUnhandled bool) {
	event := json.Get("events").GetIndex(0)
	reason := event.GetPath("severityReason", "type").MustString()
	severity := event.Get("severity").MustString()
	unhandled := event.Get("unhandled").MustBool()

	if reason != reasonType {
		t.Errorf("Wrong severity reason, expected '%s', received '%s'", reasonType, reason)
	}

	if severity != expSeverity {
		t.Errorf("Wrong severity, expected '%s', received '%s'", expSeverity, severity)
	}

	if unhandled != expUnhandled {
		t.Errorf("Wrong unhandled value, expected '%t', received '%t'", expUnhandled, unhandled)
	}
}

func generateSampleConfig(endpoint string) Configuration {
	return Configuration{
		APIKey:          testAPIKey,
		Endpoints:       Endpoints{Notify: endpoint},
		ReleaseStage:    "test",
		AppType:         "foo",
		AppVersion:      "1.2.3",
		Hostname:        "web1",
		ProjectPackages: []string{"github.com/bugsnag/bugsnag-go"},
	}
}

func crashyHandler(w http.ResponseWriter, r *http.Request) {
	c := make(chan int)
	close(c)
	c <- 1
}
