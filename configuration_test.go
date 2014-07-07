package bugsnag

import (
	"testing"
)

func TestNotifyReleaseStages(t *testing.T) {
	Configure(Configuration{ReleaseStage: "production", NotifyReleaseStages: nil})

	if !Config.notifyInReleaseStage() {
		t.Errorf("NotifyReleaseStages set to nil doesn't work")
	}

	Configure(Configuration{ReleaseStage: "production", NotifyReleaseStages: []string{"development", "production"}})
	if !Config.notifyInReleaseStage() {
		t.Errorf("NotifyReleaseStages set to array doesn't allow notifying")
	}

	Configure(Configuration{ReleaseStage: "staging", NotifyReleaseStages: []string{"development", "production"}})
	if Config.notifyInReleaseStage() {
		t.Errorf("NotifyReleaseStages set to array doesn't prevent notifying")
	}
}

func TestProjectPackages(t *testing.T) {
	Configure(Configuration{ProjectPackages: []string{"main", "github.com/ConradIrwin/*"}})
	if !Config.isProjectPackage("main") {
		t.Errorf("literal project package doesn't work")
	}
	if !Config.isProjectPackage("github.com/ConradIrwin/foo") {
		t.Errorf("wildcard project package doesn't work")
	}
	if Config.isProjectPackage("runtime") {
		t.Errorf("wrong packges being marked in project")
	}
	if Config.isProjectPackage("github.com/ConradIrwin/foo/bar") {
		t.Errorf("wrong packges being marked in project")
	}

}
