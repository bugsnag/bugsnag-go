package sessions

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	uuid "github.com/gofrs/uuid"
)

const (
	sessionEndpoint = "http://localhost:9181"
	testAPIKey      = "166f5ad3590596f9aa8d601ea89af845"
)

type testHTTPClient struct {
	reqs []*http.Request
}

// A simple io.ReadCloser that we can inject as a body of a http.Request.
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func (c *testHTTPClient) Do(r *http.Request) (*http.Response, error) {
	c.reqs = append(c.reqs, r)
	return &http.Response{Body: nopCloser{}, StatusCode: 200}, nil
}

func TestSendsCorrectPayloadForSmallConfig(t *testing.T) {
	sessions, earliestTime := makeSessions()
	testClient := testHTTPClient{}

	publisher := publisher{
		config: &SessionTrackingConfiguration{Endpoint: sessionEndpoint, Transport: http.DefaultTransport, APIKey: testAPIKey},
		client: &testClient,
	}

	err := publisher.publish(sessions)
	if err != nil {
		t.Error(err)
	}
	req := testClient.reqs[0]
	assertCorrectHeaders(t, req)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	root, err := simplejson.NewJson(body)
	if err != nil {
		t.Fatal(err)
	}

	hostname, _ := os.Hostname()
	notifier := root.Get("notifier")
	app := root.Get("app")
	device := root.Get("device")
	sessionCounts := root.Get("sessionCounts")
	testCases := []struct {
		property string
		got      interface{}
		exp      interface{}
	}{
		{property: "notifier.name", got: notifier.Get("name").MustString(), exp: "Bugsnag Go"},
		{property: "notifier.url", got: notifier.Get("url").MustString(), exp: "https://github.com/bugsnag/bugsnag-go"},
		{property: "notifier.version", got: notifier.Get("version").MustString(), exp: ""},

		{property: "app.type", got: app.Get("type").MustString(), exp: ""},
		{property: "app.releaseStage", got: app.Get("releaseStage").MustString(), exp: "production"},
		{property: "app.version", got: app.Get("version").MustString(), exp: ""},

		{property: "device.osName", got: device.Get("osName").MustString(), exp: runtime.GOOS},
		{property: "device.hostname", got: device.Get("hostname").MustString(), exp: hostname},
		{property: "sessionCounts.startedAt", got: sessionCounts.Get("startedAt").MustString(), exp: earliestTime},
		{property: "sessionCounts.sessionsStarted", got: sessionCounts.Get("sessionsStarted").MustInt(), exp: len(sessions)},
	}
	for _, tc := range testCases {
		t.Run(tc.property, func(st *testing.T) {
			if tc.got != tc.exp {
				t.Errorf("Expected property '%s' in JSON to be '%v' but was '%v'", tc.property, tc.exp, tc.got)
			}
		})
	}
}

func TestSendsCorrectPayloadForBigConfig(t *testing.T) {
	sessions, earliestTime := makeSessions()

	testClient := testHTTPClient{}
	publisher := publisher{
		config: makeHeavyConfig(),
		client: &testClient,
	}

	err := publisher.publish(sessions)
	if err != nil {
		t.Error(err)
	}
	req := testClient.reqs[0]
	assertCorrectHeaders(t, req)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	root, err := simplejson.NewJson(body)
	if err != nil {
		t.Fatal(err)
	}

	notifier := root.Get("notifier")
	device := root.Get("device")
	app := root.Get("app")
	sessionCounts := root.Get("sessionCounts")
	testCases := []struct {
		property string
		expected interface{}
		got      interface{}
	}{
		{property: "notifier.name", got: notifier.Get("name").MustString(), expected: "Bugsnag Go"},
		{property: "notifier.url", got: notifier.Get("url").MustString(), expected: "https://github.com/bugsnag/bugsnag-go"},
		{property: "notifier.version", got: notifier.Get("version").MustString(), expected: "2.3.4-alpha"},
		{property: "app.type", got: app.Get("type").MustString(), expected: "gin"},
		{property: "app.releaseStage", got: app.Get("releaseStage").MustString(), expected: "staging"},
		{property: "app.version", got: app.Get("version").MustString(), expected: "1.2.3-beta"},
		{property: "device.osName", got: device.Get("osName").MustString(), expected: runtime.GOOS},
		{property: "device.hostname", got: device.Get("hostname").MustString(), expected: "gce-1234-us-west-1"},
		{property: "sessionCounts.startedAt", got: sessionCounts.Get("startedAt").MustString(), expected: earliestTime},
		{property: "sessionCounts.sessionsStarted", got: sessionCounts.Get("sessionsStarted").MustInt(), expected: len(sessions)},
	}
	for _, tc := range testCases {
		t.Run(tc.property, func(st *testing.T) {
			if tc.got != tc.expected {
				t.Errorf("Expected property '%s' in JSON to be '%v' but was '%v'", tc.property, tc.expected, tc.got)
			}
		})
	}
}

func makeHeavyConfig() *SessionTrackingConfiguration {
	return &SessionTrackingConfiguration{
		AppType:      "gin",
		APIKey:       testAPIKey,
		AppVersion:   "1.2.3-beta",
		Version:      "2.3.4-alpha",
		Endpoint:     sessionEndpoint,
		Transport:    http.DefaultTransport,
		ReleaseStage: "staging",
		Hostname:     "gce-1234-us-west-1",
	}
}

func makeSessions() ([]*Session, string) {
	earliestTime := time.Now().Add(-6 * time.Minute)
	genUUID := func() uuid.UUID { sessionID, _ := uuid.NewV4(); return sessionID }
	return []*Session{
		{StartedAt: earliestTime, ID: genUUID()},
		{StartedAt: earliestTime.Add(2 * time.Minute), ID: genUUID()},
		{StartedAt: earliestTime.Add(4 * time.Minute), ID: genUUID()},
	}, earliestTime.UTC().Format(time.RFC3339)
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
