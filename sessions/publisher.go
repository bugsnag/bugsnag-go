package sessions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bugsnag/bugsnag-go/headers"
)

const sessionPayloadVersion = "1"

type sessionPublisher interface {
	publish(sessions []session) error
}

type defaultPublisher struct {
	config SessionTrackingConfiguration
}

func (p *defaultPublisher) publish(sessions []session) error {
	sp := makeSessionPayload(sessions, p.config)
	buf, err := json.Marshal(sp)
	if err != nil {
		return fmt.Errorf("bugsnag/sessions/publisher.publish unable to marshal json: %v", err)
	}
	client := http.Client{Transport: p.config.Transport}
	req, err := http.NewRequest("POST", p.config.Endpoint, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("bugsnag/sessions/publisher.publish unable to create request: %v", err)
	}
	for k, v := range headers.PrefixedHeaders(p.config.APIKey, sessionPayloadVersion) {
		req.Header.Add(k, v)
	}
	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("bugsnag/sessions/publisher.publish unable to deliver session: %v", err)
	}
	return nil
}
