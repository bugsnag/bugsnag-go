package bugsnag

import (
	"encoding/json"
)

// Payload is a wrapper around Event and Configuration data
type Payload struct {
	*Event
	*Configuration
}

type hash map[string]interface{}

func (p *Payload) MarshalJSON() ([]byte, error) {

	severityReason := hash{
		"type": p.handledState.SeverityReason,
	}
	if p.handledState.Framework != "" {
		severityReason["attributes"] = hash{
			"framework": p.handledState.Framework,
		}
	}

	data := hash{
		"apiKey": p.APIKey,

		"notifier": hash{
			"name":    "Bugsnag Go",
			"url":     "https://github.com/bugsnag/bugsnag-go",
			"version": VERSION,
		},

		"events": []hash{
			{
				"payloadVersion": "2",
				"exceptions": []hash{
					{
						"errorClass": p.ErrorClass,
						"message":    p.Message,
						"stacktrace": p.Stacktrace,
					},
				},
				"severity":       p.Severity.String,
				"severityReason": severityReason,
				"unhandled":      p.handledState.Unhandled,
				"app": hash{
					"releaseStage": p.ReleaseStage,
				},
				"user":     p.User,
				"metaData": p.MetaData.sanitize(p.ParamsFilters),
			},
		},
	}

	event := data["events"].([]hash)[0]

	if p.Context != "" {
		event["context"] = p.Context
	}
	if p.GroupingHash != "" {
		event["groupingHash"] = p.GroupingHash
	}
	if p.Hostname != "" {
		event["device"] = hash{
			"hostname": p.Hostname,
		}
	}
	if p.AppType != "" {
		event["app"].(hash)["type"] = p.AppType
	}
	if p.AppVersion != "" {
		event["app"].(hash)["version"] = p.AppVersion
	}
	return json.Marshal(data)

}
