package bugsnag

import "time"

func bugsnagPrefixedHeaders(apiKey, payloadVersion string) map[string]string {
	return map[string]string{
		"Content-Type":            "application/json",
		"Bugsnag-Api-Key":         apiKey,
		"Bugsnag-Payload-Version": payloadVersion,
		"Bugsnag-Sent-At":         time.Now().Format(time.RFC3339),
	}
}
