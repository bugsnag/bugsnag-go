package sessions

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	uuid "github.com/google/uuid"
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
	return &http.Response{Body: nopCloser{}, StatusCode: 202}, nil
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

	for prop, exp := range map[string]string{
		"notifier.name":    "Bugsnag Go",
		"notifier.url":     "https://github.com/bugsnag/bugsnag-go",
		"notifier.version": "",
		"app.type":         "",
		"app.releaseStage": "production",
		"app.version":      "",
		"device.osName":    runtime.GOOS,
		"device.hostname":  hostname,
	} {
		t.Run(prop, func(st *testing.T) {
			if got := getString(root, prop); got != exp {
				t.Errorf("Expected property '%s' in JSON to be '%v' but was '%v'", prop, exp, got)
			}
		})
	}
	sessionCounts := getIndex(root, "sessionCounts", 0)
	if got, exp := getString(sessionCounts, "startedAt"), earliestTime; got != exp {
		t.Errorf("Expected sessionCounts[0].startedAt to be '%s' but was '%s'", exp, got)
	}
	if got, exp := getInt(sessionCounts, "sessionsStarted"), len(sessions); got != exp {
		t.Errorf("Expected sessionCounts[0].sessionsStarted to be %d but was %d", exp, got)
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

	for prop, exp := range map[string]string{
		"notifier.name":    "Bugsnag Go",
		"notifier.url":     "https://github.com/bugsnag/bugsnag-go",
		"notifier.version": "2.3.4-alpha",
		"app.type":         "gin",
		"app.releaseStage": "development",
		"app.version":      "1.2.3-beta",
		"device.osName":    runtime.GOOS,
		"device.hostname":  "gce-1234-us-west-1",
	} {
		t.Run(prop, func(st *testing.T) {
			if got := getString(root, prop); got != exp {
				t.Errorf("Expected property '%s' in JSON to be '%v' but was '%v'", prop, exp, got)
			}
		})
	}
	sessionCounts := getIndex(root, "sessionCounts", 0)
	if got, exp := getString(sessionCounts, "startedAt"), earliestTime; got != exp {
		t.Errorf("Expected sessionCounts[0].startedAt to be '%s' but was '%s'", exp, got)
	}
	if got, exp := getInt(sessionCounts, "sessionsStarted"), len(sessions); got != exp {
		t.Errorf("Expected sessionCounts[0].sessionsStarted to be %d but was %d", exp, got)
	}
}

func TestNoSessionsSentWhenAPIKeyIsMissing(t *testing.T) {
	sessions, _ := makeSessions()
	config := makeHeavyConfig()
	config.APIKey = "labracadabrador"
	publisher := publisher{config: config, client: &testHTTPClient{}}
	if err := publisher.publish(sessions); err != nil {
		if got, exp := err.Error(), "bugsnag/sessions/publisher.publish invalid API key: 'labracadabrador'"; got != exp {
			t.Errorf(`Expected error message "%s" but got "%s"`, exp, got)
		}
	} else {
		t.Errorf("Expected error message but no errors were returned")
	}
}

func TestNoSessionsOutsideNotifyReleaseStages(t *testing.T) {
	sessions, _ := makeSessions()

	testClient := testHTTPClient{}
	config := makeHeavyConfig()
	config.NotifyReleaseStages = []string{"staging", "production"}
	publisher := publisher{
		config: config,
		client: &testClient,
	}

	err := publisher.publish(sessions)
	if err != nil {
		t.Error(err)
	}
	if got := len(testClient.reqs); got != 0 {
		t.Errorf("Didn't expect any sessions being sent as as 'development' is outside of the notify release stages, but got %d sessions", got)
	}
}

func TestReleaseStageNotSetSendsSessionsRegardlessOfNotifyReleaseStages(t *testing.T) {
	sessions, _ := makeSessions()

	testClient := testHTTPClient{}
	config := makeHeavyConfig()
	config.NotifyReleaseStages = []string{"staging", "production"}
	config.ReleaseStage = ""
	publisher := publisher{
		config: config,
		client: &testClient,
	}

	err := publisher.publish(sessions)
	if err != nil {
		t.Error(err)
	}
	if exp, got := 1, len(testClient.reqs); got != exp {
		t.Errorf("Expected %d sessions sent when the release stage is \"\" regardless of notify release stage, but got %d", exp, got)
	}
}

func makeHeavyConfig() *SessionTrackingConfiguration {
	return &SessionTrackingConfiguration{
		AppType:             "gin",
		APIKey:              testAPIKey,
		AppVersion:          "1.2.3-beta",
		Version:             "2.3.4-alpha",
		Endpoint:            sessionEndpoint,
		Transport:           http.DefaultTransport,
		ReleaseStage:        "development",
		Hostname:            "gce-1234-us-west-1",
		NotifyReleaseStages: []string{"development"},
	}
}

func makeSessions() ([]*Session, string) {
	earliestTime := time.Now().Add(-6 * time.Minute)
	return []*Session{
		{StartedAt: earliestTime, ID: uuid.New()},
		{StartedAt: earliestTime.Add(2 * time.Minute), ID: uuid.New()},
		{StartedAt: earliestTime.Add(4 * time.Minute), ID: uuid.New()},
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
