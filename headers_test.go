package bugsnag

import (
	"testing"
	"time"
)

const APIKey = "abcd1234abcd1234"

func TestConstantBugsnagPrefixedHeaders(t *testing.T) {
	headers := bugsnagPrefixedHeaders(APIKey)
	testCases := []struct {
		header   string
		expected string
	}{
		{header: "Content-Type", expected: "application/json"},
		{header: "Bugsnag-Api-Key", expected: APIKey},
		{header: "Bugsnag-Payload-Version", expected: "1"},
	}
	for _, tc := range testCases {
		t.Run(tc.header, func(st *testing.T) {
			if got := headers[tc.header]; got != tc.expected {
				t.Errorf("Expected headers to contain %s header %s but was %s", tc.header, tc.expected, got)
			}
		})
	}
}

func TestTimeDependentBugsnagPrefixedHeaders(t *testing.T) {
	headers := bugsnagPrefixedHeaders(APIKey)
	sentAtString := headers["Bugsnag-Sent-At"]
	sentAt, err := time.Parse(time.RFC3339, sentAtString)

	if err != nil {
		t.Errorf("Error when attempting to parse Bugsnag-Sent-At header: %s", sentAtString)
	}

	if now := time.Now(); now.Sub(sentAt) > time.Second || now.Sub(sentAt) < -time.Second {
		t.Errorf("Expected Bugsnag-Sent-At header approx. %s but was %s", now.UTC().Format(time.RFC3339), sentAtString)
	}
}
