package sessions_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/sessions/internal"
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
		root, err := testutil.ExtractPayload(r)
		if err != nil {
			t.Fatal(err)
		}
		hostname, _ := os.Hostname()
		testCases := []struct {
			property string
			expected string
		}{
			{property: "notifier.name", expected: "Bugsnag Go"},
			{property: "notifier.url", expected: "https://github.com/bugsnag/bugsnag-go"},
			{property: "notifier.version", expected: bugsnag.VERSION},
			{property: "app.type", expected: ""},
			{property: "app.releaseStage", expected: "production"},
			{property: "app.version", expected: ""},
			{property: "device.osName", expected: runtime.GOOS},
			{property: "device.hostname", expected: hostname},
		}
		for _, tc := range testCases {
			t.Run(tc.property, func(st *testing.T) {
				got, err := testutil.GetJSONString(root, tc.property)
				if err != nil {
					t.Error(err)
				}
				if got != tc.expected {
					t.Errorf("Expected property '%s' in JSON to be '%s' but was '%s'", tc.property, tc.expected, got)
				}
			})
		}
		mutex.Lock()
		defer mutex.Unlock()
		sessionsStarted += getSessionsStarted(t, root)
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

func getSessionsStarted(t *testing.T, root *json.RawMessage) int {
	subobj, err := testutil.GetNestedJSON(root, "sessionCounts")
	if err != nil {
		t.Error(err)
		return 0
	}
	var sessionCounts map[string]*json.RawMessage
	err = json.Unmarshal(*subobj, &sessionCounts)
	if err != nil {
		t.Error(err)
		return 0
	}
	var got int
	err = json.Unmarshal(*sessionCounts["sessionsStarted"], &got)
	if err != nil {
		t.Error(err)
	}
	return got
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
