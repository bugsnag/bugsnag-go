package bugsnag

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type payload struct {
	*Event
	*Configuration
}

type hash map[string]interface{}

func (p *payload) deliver() {

	if len(p.APIKey) != 32 {
		p.log("bugsnag/payload.deliver: invalid api key: %#v\n", p.APIKey)
		return
	}

	buf, err := json.Marshal(p)

	if err != nil {
		p.log("bugsnag/payload.deliver: %s\n", err)
		return
	}

	resp, err := http.Post(p.Endpoint, "application/json", bytes.NewBuffer(buf))

	if err != nil {
		p.log("bugsnag/payload.deliver: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		p.log("bugsnag/payload.deliver: Got HTTP %s\n", resp.Status)
		return
	}
}

func (p *payload) MarshalJSON() ([]byte, error) {

	data := hash{
		"apiKey": p.APIKey,

		"notifier": hash{
			"name":    "Bugsnag Go",
			"url":     "https://github.com/bugsnag/bugsnag-go",
			"version": "0.1",
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
				"severity": p.Severity.String,
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
	if p.AppVersion != "" {
		event["app"].(hash)["version"] = p.AppVersion
	}
	return json.Marshal(data)

}
