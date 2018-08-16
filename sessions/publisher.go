package sessions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bugsnag/bugsnag-go/headers"
)

const sessionPayloadVersion = "1.0"

type sessionPublisher interface {
	publish(sessions []session) error
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type publisher struct {
	config *SessionTrackingConfiguration
	client httpClient
}

func (p *publisher) publish(sessions []session) error {
	payload := makeSessionPayload(sessions, p.config)
	buf, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("bugsnag/sessions/publisher.publish unable to marshal json: %v", err)
	}
	req, err := http.NewRequest("POST", p.config.Endpoint, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("bugsnag/sessions/publisher.publish unable to create request: %v", err)
	}
	for k, v := range headers.PrefixedHeaders(p.config.APIKey, sessionPayloadVersion) {
		req.Header.Add(k, v)
	}
	res, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("bugsnag/sessions/publisher.publish unable to deliver session: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("bugsnag/session.deliverSessions got HTTP %s", res.Status)
	}
	return nil
}
