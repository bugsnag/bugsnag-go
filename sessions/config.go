package sessions

import (
	"net/http"
	"time"
)

// SessionTrackingConfiguration defines the configuration required for session tracking.
type SessionTrackingConfiguration struct {
	PublishInterval time.Duration

	APIKey   string
	Endpoint string
	Version  string

	ReleaseStage string
	Hostname     string
	AppType      string
	AppVersion   string
	Transport    http.RoundTripper
}

// Update modifies the values inside the struct to match the configured keys of the given config.
// Existing blank keys will not be cleared.
func (c *SessionTrackingConfiguration) Update(config *SessionTrackingConfiguration) {
	if config.PublishInterval != 0 {
		c.PublishInterval = config.PublishInterval
	}
	if config.APIKey != "" {
		c.APIKey = config.APIKey
	}
	if config.Endpoint != "" {
		c.Endpoint = config.Endpoint
	}
	if config.Version != "" {
		c.Version = config.Version
	}
	if config.ReleaseStage != "" {
		c.ReleaseStage = config.ReleaseStage
	}
	if config.Hostname != "" {
		c.Hostname = config.Hostname
	}
	if config.AppType != "" {
		c.AppType = config.AppType
	}
	if config.AppVersion != "" {
		c.AppVersion = config.AppVersion
	}
	if config.Transport != nil {
		c.Transport = config.Transport
	}
}
