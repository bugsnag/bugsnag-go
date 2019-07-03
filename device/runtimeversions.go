package device

import (
	"runtime"
)

// Cached runtime versions
var versions *RuntimeVersions

// RuntimeVersions define the various versions of Go and any framework that may
// be in use.
// As a user of the notifier you're unlikely to need to modify this struct.
// As such, the authors reserve the right to introduce breaking changes to the
// properties in this struct, in particular the framework versions are liable
// to change in new versions of the notifier, in minor/patch versions. You have
// been warned.
type RuntimeVersions struct {
	Go string `json:"go"`

	Gin     string `json:"gin,omitempty"`
	Martini string `json:"martini,omitempty"`
	Negroni string `json:"negroni,omitempty"`
	Revel   string `json:"revel,omitempty"`
}

// GetRuntimeVersions retrieves the recorded runtime versions in a goroutine-safe manner.
func GetRuntimeVersions() *RuntimeVersions {
	if versions == nil {
		versions = &RuntimeVersions{Go: runtime.Version()}
	}
	return versions
}

// AddVersion permits a framework to register its version, assuming it's one of
// the officially supported frameworks.
func AddVersion(framework, version string) {
	if versions == nil {
		versions = &RuntimeVersions{Go: runtime.Version()}
	}
	switch framework {
	case "martini":
		versions.Martini = version
	case "gin":
		versions.Gin = version
	case "negroni":
		versions.Negroni = version
	case "revel":
		versions.Revel = version
	}
}
