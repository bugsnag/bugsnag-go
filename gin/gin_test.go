package bugsnaggin_test

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

	"github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/gin"
	"github.com/gin-gonic/gin"
)

const (
	testAPIKey = "166f5ad3590596f9aa8d601ea89af845"
	port       = "9079"
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

func TestGin(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	g := gin.Default()

	userID := "1234abcd"
	g.Use(bugsnaggin.AutoNotify(bugsnag.Configuration{
		APIKey:    testAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	}, bugsnag.User{Id: userID}))

	g.GET("/unhandled", performUnhandledCrash)
	g.GET("/handled", performHandledError)
	go g.Run("localhost:9079") //This call blocks

	t.Run("AutoNotify", func(st *testing.T) {
		time.Sleep(1 * time.Second)
		_, err := http.Get("http://localhost:9079/unhandled")
		if err != nil {
			t.Error(err)
		}
		report := <-reports
		r, _ := simplejson.NewJson(report)
		hostname, _ := os.Hostname()
		assertPayload(st, r, fmt.Sprintf(`
		{
			"apiKey":"166f5ad3590596f9aa8d601ea89af845",
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
						"request":{ "httpMethod":"GET", "url":"http://localhost:9079/unhandled" }
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
		`, hostname, userID, bugsnag.VERSION))
	})

	t.Run("Manual notify", func(st *testing.T) {
		_, err := http.Get("http://localhost:9079/handled")
		if err != nil {
			t.Error(err)
		}
		report := <-reports
		r, _ := simplejson.NewJson(report)
		hostname, _ := os.Hostname()
		assertPayload(st, r, fmt.Sprintf(`
		{
			"apiKey":"166f5ad3590596f9aa8d601ea89af845",
			"events":[
				{
					"app":{ "releaseStage":"" },
					"context":"/handled",
					"device":{ "hostname": "%s" },
					"exceptions":[
						{
							"errorClass":"*errors.errorString",
							"message":"Ooopsie",
							"stacktrace":[]
						}
					],
					"payloadVersion":"4",
					"severity":"warning",
					"severityReason":{ "type":"handledError" },
					"unhandled":false,
					"user":{ "id": "%s" }
				}
			],
			"notifier":{
				"name":"Bugsnag Go",
				"url":"https://github.com/bugsnag/bugsnag-go",
				"version": "%s"
			}
		}
		`, hostname, "987zyx", bugsnag.VERSION))
	})
}

func performHandledError(c *gin.Context) {
	ctx := bugsnag.StartSession(context.Background())
	bugsnag.Notify(ctx, fmt.Errorf("Ooopsie"), bugsnag.User{Id: "987zyx"})
}

func performUnhandledCrash(c *gin.Context) {
	c.String(http.StatusOK, "OK")
	var a struct{}
	crash(a)
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

func checkFrame(t *testing.T, frame *simplejson.Json, exp stackFrame) {
	if got := getString(frame, "file"); got != exp.file {
		t.Errorf("Expected frame file to be '%s' but was '%s'", exp.file, got)
	}
	if got := getString(frame, "method"); got != exp.method {
		t.Errorf("Expected frame method to be '%s' but was '%s'", exp.method, got)
	}
	if got := getInt(frame, "lineNumber"); got != exp.lineNumber && exp.inProject { // Don't check files that vary per version of go
		t.Errorf("Expected frame line number to be %d but was %d", exp.lineNumber, got)
	}
	if got := getBool(frame, "inProject"); got != exp.inProject {
		t.Errorf("Expected frame inProject to be '%v' but was '%v'", exp.inProject, got)
	}
}
