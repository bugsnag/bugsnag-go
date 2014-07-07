package bugsnag

import (
	"testing"
)

func TestNotifyReleaseStages(t *testing.T) {
	if !(&Configuration{ReleaseStage: "production", NotifyReleaseStages: nil}).notifyInReleaseStage() {
		t.Errorf("NotifyReleaseStages set to nil doesn't work")
	}

	if !(&Configuration{ReleaseStage: "production", NotifyReleaseStages: []string{"development", "production"}}).notifyInReleaseStage() {
		t.Errorf("NotifyReleaseStages set to array doesn't allow notifying")
	}

	if (&Configuration{ReleaseStage: "development", NotifyReleaseStages: []string{"production"}}).notifyInReleaseStage() {
		t.Errorf("NotifyReleaseStages set to array doesn't prevent notifying")
	}
}

func TestProjectPackages(t *testing.T) {
	config := &Configuration{ProjectPackages: []string{"main", "github.com/ConradIrwin/*"}}
	if !config.isProjectPackage("main") {
		t.Errorf("literal project package doesn't work")
	}
	if !config.isProjectPackage("github.com/ConradIrwin/foo") {
		t.Errorf("wildcard project package doesn't work")
	}
	if config.isProjectPackage("runtime") {
		t.Errorf("wrong packges being marked in project")
	}
	if config.isProjectPackage("github.com/ConradIrwin/foo/bar") {
		t.Errorf("wrong packges being marked in project")
	}

}
