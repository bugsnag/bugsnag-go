package sessions

import (
	"os"
	"runtime"
	"time"
)

type notifierPayload struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Version string `json:"version"`
}

type appPayload struct {
	Type         string `json:"type"`
	ReleaseStage string `json:"releaseStage"`
	Version      string `json:"version"`
}

type devicePayload struct {
	OsName   string `json:"osName"`
	Hostname string `json:"hostname"`
}

type sessionCountsPayload struct {
	StartedAt       string `json:"startedAt"`
	SessionsStarted int    `json:"sessionsStarted"`
}

type sessionPayload struct {
	Notifier      notifierPayload      `json:"notifier"`
	App           appPayload           `json:"app"`
	Device        devicePayload        `json:"device"`
	SessionCounts sessionCountsPayload `json:"sessionCounts"`
}

func makeSessionPayload(sessions []session, config SessionTrackingConfiguration) sessionPayload {
	releaseStage := config.ReleaseStage
	if releaseStage == "" {
		releaseStage = "production"
	}
	hostname := config.Hostname
	if hostname == "" {
		hostname, _ = os.Hostname() //Ignore the hostname if this call errors
	}

	return sessionPayload{
		Notifier: notifierPayload{
			Name:    "Bugsnag Go",
			URL:     "https://github.com/bugsnag/bugsnag-go",
			Version: config.Version,
		},
		App: appPayload{
			Type:         config.AppType,
			Version:      config.AppVersion,
			ReleaseStage: releaseStage,
		},
		Device: devicePayload{
			OsName:   runtime.GOOS,
			Hostname: hostname,
		},
		SessionCounts: sessionCountsPayload{
			StartedAt:       sessions[0].startedAt.UTC().Format(time.RFC3339),
			SessionsStarted: len(sessions),
		},
	}
}
