package headers

import (
	"fmt"
	"time"
)

// PrefixedHeaders returns a map of Content-Type and the 'Bugsnag-' headers for
// API key, payload version, and the time at which the request is being sent.
func PrefixedHeaders(apiKey, payloadVersion, sha1 string) map[string]string {
	integrityHeader := fmt.Sprintf("sha1 %v", sha1)

	return map[string]string{
		"Content-Type":            "application/json",
		"Bugsnag-Api-Key":         apiKey,
		"Bugsnag-Payload-Version": payloadVersion,
		"Bugsnag-Sent-At":         time.Now().UTC().Format(time.RFC3339),
		"Bugsnag-Integrity":       integrityHeader,
	}
}
