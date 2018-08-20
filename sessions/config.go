package sessions

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// SessionTrackingConfiguration defines the configuration options relevant for session tracking.
// These are likely a subset of the global bugsnag.Configuration. Users should
// not modify this struct directly but rather call
// `bugsnag.Configure(bugsnag.Configuration)` which will update this configuration in return.
type SessionTrackingConfiguration struct {
	// PublishInterval defines how often the sessions are sent off to the session server.
	PublishInterval time.Duration

	// APIKey defines the API key for the Bugsnag project. Same value as for reporting errors.
	APIKey string
	// Endpoint is the URI of the session server to receive session payloads.
	Endpoint string
	// Version defines the current version of the notifier.
	Version string

	// ReleaseStage defines the release stage, e.g. "production" or "staging",
	// that this session occurred in. The release stage, in combination with
	// the app version make up the release that Bugsnag tracks.
	ReleaseStage string
	// Hostname defines the host of the server this application is running on.
	Hostname string
	// AppType defines the type of the application.
	AppType string
	// AppVersion defines the version of the application.
	AppVersion string
	// Transport defines the http.RoundTripper to be used for managing HTTP requests.
	Transport http.RoundTripper

	// Logger is the logger that Bugsnag should log to. Uses the same defaults
	// as go's builtin logging package. This logger gets invoked when any error
	// occurs inside the library itself.
	Logger interface {
		Printf(format string, v ...interface{})
	}

	mutex sync.Mutex
}

// Update modifies the values inside the receiver to match the non-default properties of the given config.
// Existing properties will not be cleared when given empty fields.
func (c *SessionTrackingConfiguration) Update(config *SessionTrackingConfiguration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
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
	if config.Logger != nil {
		c.Logger = config.Logger
	}
}

func (c *SessionTrackingConfiguration) logf(fmt string, args ...interface{}) {
	if c != nil && c.Logger != nil {
		c.Logger.Printf(fmt, args...)
	} else {
		log.Printf(fmt, args...)
	}
}
