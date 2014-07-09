package bugsnag

import (
	"log"
	"net/http"
	"strings"
)

type Configuration struct {
	// The API key, e.g. "c9d60ae4c7e70c4b6c4ebd3e8056d2b8"
	APIKey string
	// The Bugsnag endpoint, default "https://notify.bugsnag.com/"
	Endpoint string

	// The current release stage
	ReleaseStage string
	// The currently running version of the app
	AppVersion string
	// The hostname of the current server
	Hostname string

	// Release stages to notify in, default nil implies all release stages.
	NotifyReleaseStages []string
	// packages to consider in-project, default: {"main"}
	ProjectPackages []string
	// keys to filter out of meta-data, default: {"password", "secret"}
	ParamsFilters []string

	// A function to install a PanicHandler, defaults to panicwrap.
	PanicHandler func()

	// The logger to use, defaults to the global logger
	Logger *log.Logger
	// The http Transport to use, defaults to the default http Transport
	Transport http.RoundTripper
	// TODO: remember to update the update() function when modifying this struct
}

func (config *Configuration) update(other *Configuration) *Configuration {
	if other.APIKey != "" {
		config.APIKey = other.APIKey
	}
	if other.Endpoint != "" {
		config.Endpoint = other.Endpoint
	}
	if other.Hostname != "" {
		config.Hostname = other.Hostname
	}
	if other.AppVersion != "" {
		config.AppVersion = other.AppVersion
	}
	if other.ReleaseStage != "" {
		config.ReleaseStage = other.ReleaseStage
	}
	if other.ParamsFilters != nil {
		config.ParamsFilters = other.ParamsFilters
	}
	if other.ProjectPackages != nil {
		config.ProjectPackages = other.ProjectPackages
	}
	if other.Logger != nil {
		config.Logger = other.Logger
	}
	if other.NotifyReleaseStages != nil {
		config.NotifyReleaseStages = other.NotifyReleaseStages
	}
	if other.PanicHandler != nil {
		config.PanicHandler = other.PanicHandler
	}
	if other.Transport != nil {
		config.Transport = other.Transport
	}

	return config
}

func (config *Configuration) merge(other *Configuration) *Configuration {
	return config.clone().update(other)
}

func (config *Configuration) clone() *Configuration {
	clone := *config
	return &clone
}

func (config *Configuration) isProjectPackage(pkg string) bool {
	for _, p := range config.ProjectPackages {
		if p == pkg {
			return true
		} else if len(p) > 2 && p[len(p)-2] == '/' && p[len(p)-1] == '*' {
			idx := strings.LastIndex(pkg, "/")
			if idx > -1 && pkg[:idx] == p[:len(p)-2] {
				return true
			}
		}
	}
	return false
}

func (config *Configuration) stripProjectPackages(file string) string {
	for _, p := range config.ProjectPackages {
		if len(p) > 2 && p[len(p)-2] == '/' && p[len(p)-1] == '*' {
			p = p[:len(p)-1]
		} else {
			p = p + "/"
		}
		if strings.HasPrefix(file, p) {
			return strings.TrimPrefix(file, p)
		}
	}

	return file
}

func (config *Configuration) log(fmt string, args ...interface{}) {
	if config != nil && config.Logger != nil {
		config.Logger.Printf(fmt, args...)
	} else {
		log.Printf(fmt, args...)
	}
}

func (config *Configuration) notifyInReleaseStage() bool {
	if config.NotifyReleaseStages == nil {
		return true
	} else {
		for _, r := range config.NotifyReleaseStages {
			if r == config.ReleaseStage {
				return true
			}
		}
		return false
	}
}
