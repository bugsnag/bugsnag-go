// Package testutil can be .-imported to gain access to useful test functions.
package testutil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

// TestAPIKey is a fake API key that can be used for testing
const TestAPIKey = "166f5ad3590596f9aa8d601ea89af845"

// Setup sets up and returns a test event server for receiving the event payloads.
// report payloads published to the returned server's URL will be put on the returned channel
func Setup() (*httptest.Server, chan []byte) {
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

// Get travels through a JSON object and returns the specified node
func Get(j *simplejson.Json, path string) *simplejson.Json {
	return j.GetPath(strings.Split(path, ".")...)
}

// GetIndex returns the n-th element of the specified path inside the given JSON object
func GetIndex(j *simplejson.Json, path string, n int) *simplejson.Json {
	return Get(j, path).GetIndex(n)
}

func getBool(j *simplejson.Json, path string) bool {
	return Get(j, path).MustBool()
}
func getInt(j *simplejson.Json, path string) int {
	return Get(j, path).MustInt()
}
func getString(j *simplejson.Json, path string) string {
	return Get(j, path).MustString()
}
func getFirstString(j *simplejson.Json, path string) string {
	return GetIndex(j, path, 0).MustString()
}

type stackFrame struct {
	file       string
	lineNumber int
	method     string
	inProject  bool
}

// AssertPayload compares the payload that was received by the event-server to
// the expected report JSON payload
func AssertPayload(t *testing.T, report *simplejson.Json, expPretty string) {
	expReport, err := simplejson.NewJson([]byte(expPretty))
	if err != nil {
		t.Fatal(err)
	}
	expEvent := GetIndex(expReport, "events", 0)
	expException := GetIndex(expEvent, "exceptions", 0)

	event := GetIndex(report, "events", 0)
	exception := GetIndex(event, "exceptions", 0)

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
