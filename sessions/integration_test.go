package sessions_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	bugsnag "github.com/bugsnag/bugsnag-go"
)

const testAPIKey = "166f5ad3590596f9aa8d601ea89af845"
const testPublishInterval = time.Millisecond * 10
const sessionsCount = 50000

func init() {
	//Naughty injection to achieve a reasonable test duration.
	bugsnag.DefaultSessionPublishInterval = testPublishInterval
}

// Spins up a session server and checks that for every call to
// bugsnag.StartSession() a session is being recorded.
func TestStartSession(t *testing.T) {
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
		notifier := json.Get("notifier")
		app := json.Get("app")
		device := json.Get("device")
		hostname, _ := os.Hostname()
		sessionCounts := json.Get("sessionCounts")
		tt := []struct {
			prop string
			exp  interface{}
			got  interface{}
		}{
			{got: notifier.Get("name").MustString(), prop: "notifier.name", exp: "Bugsnag Go"},
			{got: notifier.Get("url").MustString(), prop: "notifier.url", exp: "https://github.com/bugsnag/bugsnag-go"},
			{got: notifier.Get("version").MustString(), prop: "notifier.version", exp: bugsnag.VERSION},
			{got: app.Get("releaseStage").MustString(), prop: "app.releaseStage", exp: "production"},
			{got: app.Get("version").MustString(), prop: "app.version", exp: ""},
			{got: device.Get("osName").MustString(), prop: "device.osName", exp: runtime.GOOS},
			{got: device.Get("hostname").MustString(), prop: "device.hostname", exp: hostname},
		}
		for _, tc := range tt {
			if tc.got != tc.exp {
				t.Errorf("Expected '%s' to be '%s' but was %s", tc.prop, tc.exp, tc.got)
			}
		}
		if got := sessionCounts.Get("startedAt").MustString(); len(got) != 20 {
			t.Errorf("Expected 'sessionCounts.startedAt' to be valid timestamp but was %s", got)
		}
		mutex.Lock()
		defer mutex.Unlock()
		sessionsStarted += sessionCounts.Get("sessionsStarted").MustInt()
	}))
	defer ts.Close()

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

	time.Sleep(testPublishInterval * 2)

	mutex.Lock()
	defer mutex.Unlock()
	if got, exp := sessionsStarted, sessionsCount; got != exp {
		t.Errorf("Expected %d sessions started, but was %d", got, exp)
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
