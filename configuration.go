package bugsnag

import (
	"fmt"
	"log"
	"strings"
)

type Configuration struct {
	// The API key, e.g. "c9d60ae4c7e70c4b6c4ebd3e8056d2b8"
	APIKey string
	// The Bugsnag endpoint, default "https://notify.bugsnag.com/"
	Endpoint string

	// The hostname of the current server
	Hostname string
	// The currently running version of the app
	AppVersion string
	// The current release stage
	ReleaseStage string

	// keys to filter out of meta-data, default: {"password", "secret"}
	ParamsFilters []string

	// directory to strip from in-project frames, default: ""
	ProjectRoot string
	// packages to consider in-project, default: {"main"}
	ProjectPackages []string

	// The logger to use, defaults to the global logger
	Logger *log.Logger
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
	if other.ProjectRoot != "" {
		config.ProjectRoot = other.ProjectRoot
	}
	if other.ProjectPackages != nil {
		config.ProjectPackages = other.ProjectPackages
	}
	if other.Logger != nil {
		config.Logger = other.Logger
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
		} else if len(p) > 2 && p[len(p) - 2] == '/' && p[len(p) - 1] == '*' {
			fmt.Println(p, pkg)
			idx := strings.LastIndex(pkg, "/")
			if idx > -1 && pkg[:idx] == p[:len(p) - 2] {
				return true
			}
		}
	}
	return false
}
