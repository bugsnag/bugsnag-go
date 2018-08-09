package bugsnag

import (
	"os"
	"runtime"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestBuildsCorrectPayloadWithMinimalConfig(t *testing.T) {
	earliestTime := time.Now()
	sessions := []session{
		{startedAt: earliestTime, id: uuid.NewV4()},
		{startedAt: earliestTime.Add(2 * time.Minute), id: uuid.NewV4()},
		{startedAt: earliestTime.Add(4 * time.Minute), id: uuid.NewV4()},
	}

	sp := makeSessionPayload(sessions, Configuration{})

	hostname, _ := os.Hostname()
	testCases := []struct {
		property string
		expected string
		got      string
	}{
		{property: "notifier name", expected: "Bugsnag Go", got: sp.Notifier.Name},
		{property: "notifier URL", expected: "https://github.com/bugsnag/bugsnag-go", got: sp.Notifier.URL},
		{property: "notifier version", expected: VERSION, got: sp.Notifier.Version},
		{property: "app type", expected: "", got: sp.App.Type},
		{property: "app release stage", expected: "production", got: sp.App.ReleaseStage},
		{property: "app version", expected: "", got: sp.App.Version},
		{property: "device OS", expected: runtime.GOOS, got: sp.Device.OsName},
		{property: "device hostname", expected: hostname, got: sp.Device.Hostname},
	}

	for _, tc := range testCases {
		t.Run(tc.property, func(st *testing.T) {
			if tc.got != tc.expected {
				t.Errorf("Expected %s '%s' but got '%s'", tc.property, tc.expected, tc.got)
			}
		})
	}

	if expected, got := earliestTime.UTC().Format(time.RFC3339), sp.SessionCounts.StartedAt; got != expected {
		t.Errorf("Expected the timestamp for sessions to be the earliest timestamp (%s), but was %s", expected, got)
	}
	if expected, got := 3, sp.SessionCounts.SessionsStarted; expected != got {
		t.Errorf("Expected the count of sessions %d, but was %d", expected, got)
	}
}

func TestBuildsCorrectPayloadFromConfig(t *testing.T) {
	config := Configuration{
		AppType:      "gin",
		AppVersion:   "1.2.3-beta",
		ReleaseStage: "staging",
		Hostname:     "gce-1234-us-west-1",
	}
	sp := makeSessionPayload([]session{{startedAt: time.Now(), id: uuid.NewV4()}}, config)

	testCases := []struct {
		property string
		expected string
		got      string
	}{
		{property: "app type", expected: config.AppType, got: sp.App.Type},
		{property: "app release stage", expected: config.ReleaseStage, got: sp.App.ReleaseStage},
		{property: "app version", expected: config.AppVersion, got: sp.App.Version},
		{property: "device hostname", expected: config.Hostname, got: sp.Device.Hostname},
	}

	for _, tc := range testCases {
		t.Run(tc.property, func(st *testing.T) {
			if tc.got != tc.expected {
				t.Errorf("Expected %s '%s' but got '%s'", tc.property, tc.expected, tc.got)
			}
		})
	}
}
