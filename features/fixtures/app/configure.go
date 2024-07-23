package main

import (
	"os"
	"strconv"
	"strings"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
)

func ConfigureBugsnag() bugsnag.Configuration {
	config := bugsnag.Configuration{
		APIKey:     os.Getenv("API_KEY"),
		AppVersion: os.Getenv("APP_VERSION"),
		AppType:    os.Getenv("APP_TYPE"),
		Hostname:   os.Getenv("HOSTNAME"),
	}

	if notifyReleaseStages := os.Getenv("NOTIFY_RELEASE_STAGES"); notifyReleaseStages != "" {
		config.NotifyReleaseStages = strings.Split(notifyReleaseStages, ",")
	}

	if releaseStage := os.Getenv("RELEASE_STAGE"); releaseStage != "" {
		config.ReleaseStage = releaseStage
	}

	if filters := os.Getenv("PARAMS_FILTERS"); filters != "" {
		config.ParamsFilters = []string{filters}
	}

	sync, err := strconv.ParseBool(os.Getenv("SYNCHRONOUS"))
	if err == nil {
		config.Synchronous = sync
	}

	acs, err := strconv.ParseBool(os.Getenv("AUTO_CAPTURE_SESSIONS"))
	if err == nil {
		config.AutoCaptureSessions = acs
	}

	config.Endpoints = bugsnag.Endpoints{
		Notify:   os.Getenv("BUGSNAG_NOTIFY_ENDPOINT"),
		Sessions: os.Getenv("BUGSNAG_SESSIONS_ENDPOINT"),
	}
	
	return config
}
