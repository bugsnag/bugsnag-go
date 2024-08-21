package headers

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

const APIKey = "abcd1234abcd1234"
const testPayloadVersion = "3"
const testSHA = "5e13ae4640ae4ae0e09c05b7bb060f544dabd042"

func TestConstantBugsnagPrefixedHeaders(t *testing.T) {
	headers := PrefixedHeaders(APIKey, testPayloadVersion, testSHA)
	testCases := []struct {
		header   string
		expected string
	}{
		{header: "Content-Type", expected: "application/json"},
		{header: "Bugsnag-Api-Key", expected: APIKey},
		{header: "Bugsnag-Payload-Version", expected: testPayloadVersion},
		{header: "Bugsnag-Integrity", expected: fmt.Sprintf("sha1 %v", testSHA)},
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
	headers := PrefixedHeaders(APIKey, testPayloadVersion, testSHA)
	sentAtString := headers["Bugsnag-Sent-At"]
	if !strings.HasSuffix(sentAtString, "Z") {
		t.Errorf("Error when setting Bugsnag-Sent-At header: %s, doesn't end with a Z", sentAtString)
	}
	sentAt, err := time.Parse(time.RFC3339, sentAtString)

	if err != nil {
		t.Errorf("Error when attempting to parse Bugsnag-Sent-At header: %s", sentAtString)
	}

	if now := time.Now(); now.Sub(sentAt) > time.Second || now.Sub(sentAt) < -time.Second {
		t.Errorf("Expected Bugsnag-Sent-At header approx. %s but was %s", now.UTC().Format(time.RFC3339), sentAtString)
	}
}
