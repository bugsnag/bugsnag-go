package bugsnagmartini_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/martini"
	"github.com/go-martini/martini"
)

const (
	testAPIKey = "166f5ad3590596f9aa8d601ea89af845"
)

// setup sets up and returns a test event server for receiving the event payloads.
// report payloads published to the returned server's URL will be put on the returned channel
func setup() (*httptest.Server, chan []byte) {
	reports := make(chan []byte, 10)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Bugsnag called")
		if strings.Contains(r.URL.Path, "sessions") {
			return
		}
		body, _ := ioutil.ReadAll(r.Body)
		reports <- body
	})), reports
}

func performHandledError(notifier *bugsnag.Notifier) {
	ctx := bugsnag.StartSession(context.Background())
	notifier.Notify(ctx, fmt.Errorf("Ooopsie"), bugsnag.User{Id: "987zyx"})
}

func performUnhandledCrash() string {
	var a struct{}
	crash(a)
	return "ok"
}

func TestMartini(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	m := martini.Classic()

	userID := "1234abcd"
	m.Use(martini.Recovery())
	config := bugsnag.Configuration{
		APIKey:    testAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	}
	bugsnag.Configure(config)
	m.Use(bugsnagmartini.AutoNotify(bugsnag.User{Id: userID}))

	m.Get("/unhandled", performUnhandledCrash)
	m.Get("/handled", performHandledError)
	go m.RunOnAddr(":9077")

	t.Run("AutoNotify", func(st *testing.T) {
		time.Sleep(1 * time.Second)
		_, err := http.Get("http://localhost:9077/unhandled")
		if err != nil {
			t.Error(err)
		}
		report := <-reports
		r, _ := simplejson.NewJson(report)
		hostname, _ := os.Hostname()
		assertPayload(st, r, fmt.Sprintf(`
		{
			"apiKey": "%s",
			"events":[
				{
					"app":{ "releaseStage":"" },
					"context":"/unhandled",
					"device":{ "hostname": "%s" },
					"exceptions":[
						{
							"errorClass":"*runtime.TypeAssertionError",
							"message":"interface conversion: interface {} is struct {}, not string",
							"stacktrace":[]
						}
					],
					"metaData":{
						"request":{ "httpMethod":"GET", "url":"http://localhost:9077/unhandled" }
					},
					"payloadVersion":"4",
					"severity":"error",
					"severityReason":{ "type":"unhandledErrorMiddleware" },
					"unhandled":true,
					"user":{ "id": "%s" }
				}
			],
			"notifier":{
				"name":"Bugsnag Go",
				"url":"https://github.com/bugsnag/bugsnag-go",
				"version": "%s"
			}
		}
		`, testAPIKey, hostname, userID, bugsnag.VERSION))
	})

	t.Run("Notify", func(st *testing.T) {
		time.Sleep(1 * time.Second)
		_, err := http.Get("http://localhost:9077/handled")
		if err != nil {
			t.Error(err)
		}
		report := <-reports
		r, _ := simplejson.NewJson(report)
		hostname, _ := os.Hostname()
		assertPayload(st, r, fmt.Sprintf(`
		{
			"apiKey": "%s",
			"events":[
				{
					"app":{ "releaseStage":"" },
					"device":{ "hostname": "%s" },
					"exceptions":[
						{
							"errorClass":"*errors.errorString",
							"message":"Ooopsie",
							"stacktrace":[]
						}
					],
					"metaData":{
						"request":{ "httpMethod":"GET", "url":"http://localhost:9077/handled" }
					},
					"payloadVersion":"4",
					"severity":"error",
					"severityReason":{ "type":"unhandledErrorMiddleware" },
					"unhandled":true,
					"user":{ "id": "%s" }
				}
			],
			"notifier":{
				"name":"Bugsnag Go",
				"url":"https://github.com/bugsnag/bugsnag-go",
				"version": "%s"
			}
		}
		`, testAPIKey, hostname, "987zyx", bugsnag.VERSION))
	})

}

func main() {
	if os.Getenv("BUGSNAG_TEST_VARIANT") == "beforenotify" {
		bugsnag.OnBeforeNotify(func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			event.Severity = bugsnag.SeverityInfo
			return nil
		})
	}
	m := martini.Classic()
	m.Get("/", func() string {
		var a struct{}
		crash(a)
		return "Hello world!"
	})
	m.Use(martini.Recovery())
	m.Use(bugsnagmartini.AutoNotify(bugsnag.Configuration{
		APIKey:    "166f5ad3590596f9aa8d601ea89af845",
		Endpoints: bugsnag.Endpoints{Notify: os.Getenv("BUGSNAG_NOTIFY_ENDPOINT"), Sessions: os.Getenv("BUGSNAG_SESSIONS_ENDPOINT")},
	}))
	m.Run()
}

func crash(a interface{}) string {
	return a.(string)
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

type stackFrame struct {
	file       string
	lineNumber int
	method     string
	inProject  bool
}

// assertPayload compares the payload that was received by the event-server to
// the expected report JSON payload
func assertPayload(t *testing.T, report *simplejson.Json, expPretty string) {
	expReport, err := simplejson.NewJson([]byte(expPretty))
	if err != nil {
		t.Fatal(err)
	}
	expEvent := getIndex(expReport, "events", 0)
	expException := getIndex(expEvent, "exceptions", 0)

	event := getIndex(report, "events", 0)
	exception := getIndex(event, "exceptions", 0)

	if exp, got := getBool(expEvent, "unhandled"), getBool(event, "unhandled"); got != exp {
		t.Errorf("expected 'unhandled' to be '%v' but got '%v'", exp, got)
	}
	for _, tc := range []struct {
		prop     string
		got, exp *simplejson.Json
	}{
		{got: report, exp: expReport, prop: "apiKey"},
		{got: report, exp: expReport, prop: "notifier.name"},
		{got: report, exp: expReport, prop: "notifier.version"},
		{got: report, exp: expReport, prop: "notifier.url"},
		{got: exception, exp: expException, prop: "message"},
		{got: exception, exp: expException, prop: "errorClass"},
		{got: event, exp: expEvent, prop: "user.id"},
		{got: event, exp: expEvent, prop: "severity"},
		{got: event, exp: expEvent, prop: "severityReason.type"},
		{got: event, exp: expEvent, prop: "metaData.request.httpMethod"},
		{got: event, exp: expEvent, prop: "metaData.request.url"},
	} {
		if got, exp := getString(tc.got, tc.prop), getString(tc.exp, tc.prop); got != exp {
			t.Errorf("expected '%s' to be '%s' but was '%s'", tc.prop, exp, got)
		}
	}
	assertValidSession(t, event, getBool(expEvent, "unhandled"))
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
