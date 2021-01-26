package sessions_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	bugsnag "github.com/bugsnag/bugsnag-go/v2"
)

const testAPIKey = "166f5ad3590596f9aa8d601ea89af845"
const testPublishInterval = time.Millisecond * 200
const sessionsCount = 50000

func init() {
	//Naughty injection to achieve a reasonable test duration.
	bugsnag.DefaultSessionPublishInterval = testPublishInterval
}

func get(j *simplejson.Json, path string) *simplejson.Json {
	return j.GetPath(strings.Split(path, ".")...)
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

// Spins up a session server and checks that for every call to
// bugsnag.StartSession() a session is being recorded.
func TestStartSession(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not compatible with windows builds")
		return
	}
	sessionsStarted := 0
	mutex := sync.Mutex{}

	// Test server does all the checking of individual requests
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCorrectHeaders(t, r)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		json, err := simplejson.NewJson(body)
		if err != nil {
			t.Error(err)
		}
		hostname, _ := os.Hostname()
		tt := []struct {
			prop string
			exp  interface{}
		}{
			{prop: "notifier.name", exp: "Bugsnag Go"},
			{prop: "notifier.url", exp: "https://github.com/bugsnag/bugsnag-go"},
			{prop: "notifier.version", exp: bugsnag.Version},
			{prop: "app.releaseStage", exp: "production"},
			{prop: "app.version", exp: ""},
			{prop: "device.osName", exp: runtime.GOOS},
			{prop: "device.hostname", exp: hostname},
			{prop: "device.runtimeVersions.go", exp: runtime.Version()},
			{prop: "device.runtimeVersions.gin", exp: ""},
			{prop: "device.runtimeVersions.martini", exp: ""},
			{prop: "device.runtimeVersions.negroni", exp: ""},
			{prop: "device.runtimeVersions.revel", exp: ""},
		}
		for _, tc := range tt {
			got := getString(json, tc.prop)
			if got != tc.exp {
				t.Errorf("Expected '%s' to be '%s' but was '%s'", tc.prop, tc.exp, got)
			}
		}
		sessionCounts := getIndex(json, "sessionCounts", 0)
		if got := getString(sessionCounts, "startedAt"); len(got) != 20 {
			t.Errorf("Expected 'sessionCounts.startedAt' to be valid timestamp but was %s", got)
		}
		mutex.Lock()
		defer mutex.Unlock()
		sessionsStarted += getInt(sessionCounts, "sessionsStarted")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	time.Sleep(testPublishInterval * 2) //Allow server to start

	// Minimal config. API is mandatory, URLs point to the test server
	bugsnag.Configure(bugsnag.Configuration{
		APIKey: testAPIKey,
		Endpoints: bugsnag.Endpoints{
			Sessions: ts.URL,
			Notify:   ts.URL,
		},
	})
	for i := 0; i < sessionsCount; i++ {
		bugsnag.StartSession(context.Background())
	}

	time.Sleep(testPublishInterval * 2) //Allow all messages to be processed

	mutex.Lock()
	defer mutex.Unlock()
	// Don't expect an additional session from startup as the test server URL
	// would be different between processes
	if got, exp := sessionsStarted, sessionsCount; got != exp {
		t.Errorf("Expected %d sessions started, but was %d", exp, got)
	}
}

func assertCorrectHeaders(t *testing.T, req *http.Request) {
	testCases := []struct{ name, expected string }{
		{name: "Bugsnag-Payload-Version", expected: "1.0"},
		{name: "Content-Type", expected: "application/json"},
		{name: "Bugsnag-Api-Key", expected: testAPIKey},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(st *testing.T) {
			if got := req.Header[tc.name][0]; tc.expected != got {
				t.Errorf("Expected header '%s' to be '%s' but was '%s'", tc.name, tc.expected, got)
			}
		})
	}
	name := "Bugsnag-Sent-At"
	if req.Header[name][0] == "" {
		t.Errorf("Expected header '%s' to be non-empty but was empty", name)
	}
}
