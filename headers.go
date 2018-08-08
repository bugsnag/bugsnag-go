package bugsnag

import "time"

func bugsnagPrefixedHeaders(apiKey string) map[string]string {
	return map[string]string{
		"Content-Type":            "application/json",
		"Bugsnag-Api-Key":         apiKey,
		"Bugsnag-Payload-Version": "1",
		"Bugsnag-Sent-At":         time.Now().Format(time.RFC3339),
	}
}
