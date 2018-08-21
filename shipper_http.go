package bugsnag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// NewHTTPShipper creates a shipper to send data over http
func NewHTTPShipper() Shipper {
	return &httpShipper{}
}

type httpShipper struct {
}

func (s *httpShipper) Deliver(p *payload) error {
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

	resp, err := client.Post(
		p.Endpoint,
		"application/json",
		bytes.NewBuffer(buf),
	)

	if err != nil {
		return fmt.Errorf("bugsnag/payload.deliver: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("bugsnag/payload.deliver: Got HTTP %s\n", resp.Status)
	}

	return nil
}
