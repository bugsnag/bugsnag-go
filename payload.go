package bugsnag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bugsnag/bugsnag-go/headers"
)

const notifyPayloadVersion = "2"

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
	for k, v := range headers.PrefixedHeaders(p.APIKey, notifyPayloadVersion) {
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
	return json.Marshal(reportJSON{
		APIKey: p.APIKey,
		Events: []eventJSON{
			eventJSON{
				App: &appJSON{
					ReleaseStage: p.ReleaseStage,
					Type:         p.AppType,
					Version:      p.AppVersion,
				},
				Context: p.Context,
				Device:  &deviceJSON{Hostname: p.Hostname},
				Exceptions: []exceptionJSON{
					exceptionJSON{
						ErrorClass: p.ErrorClass,
						Message:    p.Message,
						Stacktrace: p.Stacktrace,
					},
				},
				GroupingHash:   p.GroupingHash,
				Metadata:       p.MetaData.sanitize(p.ParamsFilters),
				PayloadVersion: notifyPayloadVersion,
				Severity:       p.Severity.String,
				SeverityReason: p.severityReasonPayload(),
				Unhandled:      p.handledState.Unhandled,
				User:           p.User,
			},
		},
		Notifier: notifierJSON{
			Name:    "Bugsnag Go",
			URL:     "https://github.com/bugsnag/bugsnag-go",
			Version: VERSION,
		},
	})
}

func (p *payload) severityReasonPayload() *severityReasonJSON {
	reason, framework := p.handledState.SeverityReason, p.handledState.Framework
	if reason == "" && framework == "" {
		return nil
	}
	json := severityReasonJSON{}

	if reason != "" {
		json.Type = reason
	}
	if framework != "" {
		json.Attributes = &severityAttributesJSON{Framework: framework}
	}
	return &json
}
