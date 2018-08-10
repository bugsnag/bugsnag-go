package bugsnag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type payload struct {
	*Event
	*Configuration
}

type hash map[string]interface{}

func (p *payload) deliver() error {

	if len(p.APIKey) != 32 {
		return fmt.Errorf("bugsnag/payload.deliver: invalid api key")
	}

	buf, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("bugsnag/payload.deliver: %v", err)
	}

	client := http.Client{
		Transport: p.Transport,
	}
	req, err := http.NewRequest("POST", p.Endpoints.Notify, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("bugsnag/payload.deliver unable to create request: %v", err)
	}
	for k, v := range bugsnagPrefixedHeaders(p.APIKey) {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("bugsnag/payload.deliver: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("bugsnag/payload.deliver: Got HTTP %s", resp.Status)
	}

	return nil
}

func (p *payload) MarshalJSON() ([]byte, error) {

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
