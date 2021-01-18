package sessions

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestConfigDoesNotChangeGivenBlankValues(t *testing.T) {
	c := testConfig()
	exp := testConfig()
	c.Update(&SessionTrackingConfiguration{})
	tt := []struct {
		name     string
		expected interface{}
		got      interface{}
	}{
		{"PublishInterval", exp.PublishInterval, c.PublishInterval},
		{"APIKey", exp.APIKey, c.APIKey},
		{"Endpoint", exp.Endpoint, c.Endpoint},
		{"Version", exp.Version, c.Version},
		{"ReleaseStage", exp.ReleaseStage, c.ReleaseStage},
		{"Hostname", exp.Hostname, c.Hostname},
		{"AppType", exp.AppType, c.AppType},
		{"AppVersion", exp.AppVersion, c.AppVersion},
		{"Transport", exp.Transport, c.Transport},
		{"NotifyReleaseStages", exp.NotifyReleaseStages, c.NotifyReleaseStages},
	}
	for _, tc := range tt {
		if !reflect.DeepEqual(tc.got, tc.expected) {
			t.Errorf("Expected '%s' to be '%v' but was '%v'", tc.name, tc.expected, tc.got)
		}
	}
}

func TestConfigUpdatesGivenNonDefaultValues(t *testing.T) {
	c := testConfig()
	exp := SessionTrackingConfiguration{
		PublishInterval:     40 * time.Second,
		APIKey:              "api234",
		Endpoint:            "https://docs.bugsnag.com/platforms/go/",
		Version:             "2.7.3",
		ReleaseStage:        "Production",
		Hostname:            "Brian's Surface",
		AppType:             "Revel API",
		AppVersion:          "6.3.9",
		NotifyReleaseStages: []string{"staging", "production"},
	}
	c.Update(&exp)
	tt := []struct {
		name     string
		expected interface{}
		got      interface{}
	}{
		{"PublishInterval", exp.PublishInterval, c.PublishInterval},
		{"APIKey", exp.APIKey, c.APIKey},
		{"Endpoint", exp.Endpoint, c.Endpoint},
		{"Version", exp.Version, c.Version},
		{"ReleaseStage", exp.ReleaseStage, c.ReleaseStage},
		{"Hostname", exp.Hostname, c.Hostname},
		{"AppType", exp.AppType, c.AppType},
		{"AppVersion", exp.AppVersion, c.AppVersion},
		{"NotifyReleaseStages", exp.NotifyReleaseStages, c.NotifyReleaseStages},
	}
	for _, tc := range tt {
		if !reflect.DeepEqual(tc.got, tc.expected) {
			t.Errorf("Expected '%s' to be '%v' but was '%v'", tc.name, tc.expected, tc.got)
		}
	}
}

func testConfig() SessionTrackingConfiguration {
	return SessionTrackingConfiguration{
		PublishInterval:     20 * time.Second,
		APIKey:              "api123",
		Endpoint:            "https://bugsnag.com/jobs", //If you like what you see... ;)
		Version:             "1.6.2",
		ReleaseStage:        "Staging",
		Hostname:            "Russ's MacbookPro",
		AppType:             "Gin API",
		AppVersion:          "5.2.8",
		NotifyReleaseStages: []string{"staging", "production"},
		Transport:           http.DefaultTransport,
	}
}
