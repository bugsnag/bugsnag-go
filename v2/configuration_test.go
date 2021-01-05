package bugsnag

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestNotifyReleaseStages(t *testing.T) {

	notify := " "

	var tt = []struct {
		releaseStage        string
		notifyReleaseStages []string
		expected            bool
	}{
		{
			releaseStage: "production",
			expected:     true,
		},
		{
			releaseStage:        "production",
			notifyReleaseStages: []string{"development", "production"},
			expected:            true,
		},
		{
			releaseStage:        "staging",
			notifyReleaseStages: []string{"development", "production"},
			expected:            false,
		},
		{
			notifyReleaseStages: []string{"development", "production"},
			expected:            true,
		},
	}

	for _, tc := range tt {
		rs, nrs, exp := tc.releaseStage, tc.notifyReleaseStages, tc.expected
		config := &Configuration{ReleaseStage: rs, NotifyReleaseStages: nrs}
		if config.notifyInReleaseStage() != exp {
			if !exp {
				notify = " not "
			}
			t.Errorf("expected%sto notify when release stage is '%s' and notify release stages are '%+v'", notify, rs, nrs)
		}
	}
}

func TestIsProjectPackage(t *testing.T) {

	Configure(Configuration{ProjectPackages: []string{
		"main",
		"star*",
		"example.com/a",
		"example.com/b/*",
		"example.com/c/*/*",
		"example.com/d/**",
		"example.com/e",
	}})

	var testCases = []struct {
		Path     string
		Included bool
	}{
		{"", false},
		{"main", true},
		{"runtime", false},

		{"star", true},
		{"sta", false},
		{"starred", true},
		{"star/foo", false},

		{"example.com/a", true},

		{"example.com/b", false},
		{"example.com/b/", true},
		{"example.com/b/foo", true},
		{"example.com/b/foo/bar", false},

		{"example.com/c/foo/bar", true},
		{"example.com/c/foo/bar/baz", false},

		{"example.com/d/foo/bar", true},
		{"example.com/d/foo/bar/baz", true},

		{"example.com/e", true},
	}

	for _, s := range testCases {
		if Config.isProjectPackage(s.Path) != s.Included {
			t.Error("literal project package doesn't work:", s.Path, s.Included)
		}
	}
}

func TestStripProjectPackage(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	Configure(Configuration{
		ProjectPackages: []string{
			"main",
			"star*",
			"example.com/a",
			"example.com/b/*",
			"example.com/c/**",
		},
		SourceRoot: gopath + "/src/",
	})

	var testCases = []struct {
		File     string
		Stripped string
	}{
		{"main.go", "main.go"},
		{"runtime.go", "runtime.go"},
		{"star.go", "star.go"},

		{"example.com/a/foo.go", "foo.go"},

		{"example.com/b/foo/bar.go", "foo/bar.go"},
		{"example.com/b/foo.go", "foo.go"},

		{"example.com/x/a/b/foo.go", "example.com/x/a/b/foo.go"},

		{"example.com/c/a/b/foo.go", "a/b/foo.go"},

		{gopath + "/src/runtime.go", "runtime.go"},
		{gopath + "/src/example.com/a/foo.go", "foo.go"},
		{gopath + "/src/example.com/x/a/b/foo.go", "example.com/x/a/b/foo.go"},
		{gopath + "/src/example.com/c/a/b/foo.go", "a/b/foo.go"},
	}

	for _, tc := range testCases {
		if s := Config.stripProjectPackages(tc.File); s != tc.Stripped {
			t.Error("stripProjectPackage did not remove expected path:", tc.File, tc.Stripped, "was:", s)
		}
	}
}

func TestStripCustomSourceRoot(t *testing.T) {
	Configure(Configuration{
		ProjectPackages: []string{
			"main",
			"star*",
			"example.com/a",
			"example.com/b/*",
			"example.com/c/**",
		},
		SourceRoot: "/Users/bob/code/go/src/",
	})
	var testCases = []struct {
		File     string
		Stripped string
	}{
		{"main.go", "main.go"},
		{"runtime.go", "runtime.go"},
		{"star.go", "star.go"},

		{"example.com/a/foo.go", "foo.go"},

		{"example.com/b/foo/bar.go", "foo/bar.go"},
		{"example.com/b/foo.go", "foo.go"},

		{"example.com/x/a/b/foo.go", "example.com/x/a/b/foo.go"},

		{"example.com/c/a/b/foo.go", "a/b/foo.go"},

		{"/Users/bob/code/go/src/runtime.go", "runtime.go"},
		{"/Users/bob/code/go/src/example.com/a/foo.go", "foo.go"},
		{"/Users/bob/code/go/src/example.com/x/a/b/foo.go", "example.com/x/a/b/foo.go"},
		{"/Users/bob/code/go/src/example.com/c/a/b/foo.go", "a/b/foo.go"},
	}

	for _, tc := range testCases {
		if s := Config.stripProjectPackages(tc.File); s != tc.Stripped {
			t.Error("stripProjectPackage did not remove expected path:", tc.File, tc.Stripped, "was:", s)
		}
	}
}

type CustomTestLogger struct {
	loggedMessages []string
}

func (logger *CustomTestLogger) Printf(format string, v ...interface{}) {
	logger.loggedMessages = append(logger.loggedMessages, format)
}

func TestConfiguringCustomLogger(t *testing.T) {

	l1 := log.New(os.Stdout, "", log.Lshortfile)

	l2 := &CustomTestLogger{}

	var testCases = []struct {
		config Configuration
		notify bool
		msg    string
	}{
		{
			config: Configuration{ReleaseStage: "production", NotifyReleaseStages: []string{"development", "production"}, Logger: l1},
		},
		{
			config: Configuration{ReleaseStage: "production", NotifyReleaseStages: []string{"development", "production"}, Logger: l2},
		},
	}

	for _, testCase := range testCases {
		Configure(testCase.config)

		// call printf just to illustrate it is present as the compiler does most of the hard work
		testCase.config.Logger.Printf("hello %s", "bugsnag")

	}
}

func TestEndpointDeprecationWarning(t *testing.T) {
	defaultNotify := "https://notify.bugsnag.com/"
	defaultSessions := "https://sessions.bugsnag.com/"
	setUp := func() (*Configuration, *CustomTestLogger) {
		logger := &CustomTestLogger{}
		return &Configuration{
			Endpoints: Endpoints{
				Notify:   defaultNotify,
				Sessions: defaultSessions,
			},
			Logger: logger,
		}, logger
	}

	t.Run("Setting Endpoint gives deprecation warning", func(st *testing.T) {
		c, logger := setUp()
		config := Configuration{Endpoint: "https://endpoint.whatever.com/"}
		c.update(&config)
		if got := logger.loggedMessages; len(got) != 1 {
			st.Errorf("Expected exactly one logged message but got %d: %v", len(got), got)
		}
		got := logger.loggedMessages[0]
		for _, exp := range []string{"WARNING", "Bugsnag", "Endpoint", "Endpoints", "deprecated"} {
			if !strings.Contains(got, exp) {
				st.Errorf("Expected logger message containing '%s' when configuring but got %s.", exp, got)
			}
		}
		if got, exp := c.Endpoints.Notify, config.Endpoint; got != exp {
			st.Errorf("Expected notify endpoint '%s' but got '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, ""; got != exp {
			st.Errorf("Expected sessions endpoint '%s' but got '%s'", exp, got)
		}
	})

	t.Run("Setting Endpoints.Notify without setting Endpoints.Sessions gives session disabled warning", func(st *testing.T) {
		c, logger := setUp()
		config := Configuration{
			Endpoints: Endpoints{
				Notify: "https://notify.whatever.com/",
			},
		}
		keywords := []string{"WARNING", "Bugsnag", "notify", "No sessions"}
		c.update(&config)
		if got := len(logger.loggedMessages); got != 1 {
			st.Errorf("Expected exactly one logged message but got %d", got)
		}
		got := logger.loggedMessages[0]
		for _, exp := range keywords {
			if !strings.Contains(got, exp) {
				st.Errorf("Expected logger message containing '%s' when configuring but got %s.", exp, got)
			}
		}
		if got, exp := c.Endpoints.Notify, config.Endpoints.Notify; got != exp {
			st.Errorf("Expected notify endpoint to be '%s' but was '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, ""; got != exp {
			st.Errorf("Expected sessions endpoint to be '%s' but was '%s'", exp, got)
		}
	})

	t.Run("Setting Endpoints.Sessions without setting Endpoints.Notify should panic", func(st *testing.T) {
		c, _ := setUp()
		defer func() {
			if err := recover(); err != nil {
				got := err.(string)
				for _, exp := range []string{"FATAL", "Bugsnag", "notify", "sessions"} {
					if !strings.Contains(got, exp) {
						st.Errorf("Expected panic error containing '%s' when configuring but got %s.", exp, got)
					}
				}
			} else {
				st.Errorf("Expected a panic to happen but didn't")
			}
		}()
		c.update(&Configuration{
			Endpoints: Endpoints{
				Sessions: "https://sessions.whatever.com/",
			},
		})
	})

	t.Run("Should not complain if both Endpoints.Notify and Endpoints.Sessions are configured", func(st *testing.T) {
		notifyEndpoint, sessionsEndpoint := "https://notify.whatever.com", "https://sessions.whatever.com"
		config := Configuration{
			Endpoints: Endpoints{
				Notify:   notifyEndpoint,
				Sessions: sessionsEndpoint,
			},
		}
		c, logger := setUp()
		c.update(&config)
		if len(logger.loggedMessages) != 0 {
			st.Errorf("Did not expect any messages to be logged but logged: %v", logger.loggedMessages)
		}
		if got, exp := c.Endpoints.Notify, notifyEndpoint; got != exp {
			st.Errorf("Expected Notify endpoint: '%s', but was: '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, sessionsEndpoint; got != exp {
			st.Errorf("Expected Sessions endpoint: '%s', but was: '%s'", exp, got)
		}
	})

	t.Run("Should not complain if Endpoints are not configured", func(st *testing.T) {
		c, logger := setUp()
		c.update(&Configuration{})
		if len(logger.loggedMessages) != 0 {
			st.Errorf("Did not expect any messages to be logged but logged: %v", logger.loggedMessages)
		}
		if got, exp := c.Endpoints.Notify, defaultNotify; got != exp {
			st.Errorf("Expected Notify endpoint: '%s', but was: '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, defaultSessions; got != exp {
			st.Errorf("Expected Sessions endpoint: '%s', but was: '%s'", exp, got)
		}
	})
}

func TestIsAutoCaptureSessions(t *testing.T) {
	defaultConfig := Configuration{}
	if !defaultConfig.IsAutoCaptureSessions() {
		t.Errorf("Expected automatic session tracking to be enabled by default, but was disabled")
	}

	enabledConfig := Configuration{AutoCaptureSessions: true}
	if !enabledConfig.IsAutoCaptureSessions() {
		t.Errorf("Expected automatic session tracking to be enabled when so configured, but was disabled")
	}

	disabledConfig := Configuration{AutoCaptureSessions: false}
	if disabledConfig.IsAutoCaptureSessions() {
		t.Errorf("Expected automatic session tracking to be disabled when so configured, but enabled")
	}
}
