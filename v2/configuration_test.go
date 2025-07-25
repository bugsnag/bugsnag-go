package bugsnag

import (
	"log"
	"os"
	"runtime"
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
		config := &Configuration{ReleaseStage: rs, EnabledReleaseStages: nrs}
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

	// on windows, source lines always use '/' but GOPATH may use '\' depending
	// on user settings.
	adjustedGopath := strings.Replace(gopath, "\\", "/", -1)
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

		{adjustedGopath + "/src/runtime.go", "runtime.go"},
		{adjustedGopath + "/src/example.com/a/foo.go", "foo.go"},
		{adjustedGopath + "/src/example.com/x/a/b/foo.go", "example.com/x/a/b/foo.go"},
		{adjustedGopath + "/src/example.com/c/a/b/foo.go", "a/b/foo.go"},
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

func TestStripCustomWindowsSourceRoot(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("not compatible with non-windows builds")
		return
	}
	Configure(Configuration{
		ProjectPackages: []string{
			"main",
			"star*",
			"example.com/a",
			"example.com\\b\\*",
			"example.com/c/**",
		},
		SourceRoot: "C:\\Users\\bob\\code\\go\\src\\",
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

		{"C:/Users/bob/code/go/src/runtime.go", "runtime.go"},
		{"C:/Users/bob/code/go/src/example.com/a/foo.go", "foo.go"},
		{"C:/Users/bob/code/go/src/example.com/x/a/b/foo.go", "example.com/x/a/b/foo.go"},
		{"C:/Users/bob/code/go/src/example.com/c/a/b/foo.go", "a/b/foo.go"},
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
			config: Configuration{ReleaseStage: "production", EnabledReleaseStages: []string{"development", "production"}, Logger: l1},
		},
		{
			config: Configuration{ReleaseStage: "production", EnabledReleaseStages: []string{"development", "production"}, Logger: l2},
		},
	}

	for _, testCase := range testCases {
		Configure(testCase.config)

		// call printf just to illustrate it is present as the compiler does most of the hard work
		testCase.config.Logger.Printf("hello %s", "bugsnag")

	}
}

func TestEndpointDeprecationWarning(t *testing.T) {
	defaultNotify := "https://notify.bugsnag.com"
	defaultSessions := "https://sessions.bugsnag.com"
	setUp := func() (*Configuration, *CustomTestLogger) {
		logger := &CustomTestLogger{}
		return &Configuration{
			Logger: logger,
		}, logger
	}

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
func TestEndpointFromEnvironment(t *testing.T) {
	notifyEndpoint, sessionsEndpoint := "https://notify.custom.com", "https://sessions.custom.com"
	setUp := func() {
		os.Setenv("BUGSNAG_NOTIFY_ENDPOINT", notifyEndpoint)
		os.Setenv("BUGSNAG_SESSIONS_ENDPOINT", sessionsEndpoint)
		os.Setenv("BUGSNAG_API_KEY", "00000abcdef0123456789abcdef012345")
	}
	cleanup := func() {
		defer os.Unsetenv("BUGSNAG_NOTIFY_ENDPOINT")
		defer os.Unsetenv("BUGSNAG_SESSIONS_ENDPOINT")
		defer os.Unsetenv("BUGSNAG_API_KEY")
	}

	t.Run("Should not override endpoints set by environment variables", func(st *testing.T) {
		setUp()
		defer cleanup()
		c := &Configuration{}
		c.loadEnv()

		if got, exp := c.Endpoints.Notify, notifyEndpoint; got != exp {
			st.Errorf("Expected Notify endpoint: '%s', but was: '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, sessionsEndpoint; got != exp {
			st.Errorf("Expected Sessions endpoint: '%s', but was: '%s'", exp, got)
		}

		c.update(&Configuration{
			ProjectPackages: []string{"main"},
		})
		if got, exp := c.Endpoints.Notify, notifyEndpoint; got != exp {
			st.Errorf("Expected Notify endpoint: '%s', but was: '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, sessionsEndpoint; got != exp {
			st.Errorf("Expected Sessions endpoint: '%s', but was: '%s'", exp, got)
		}
	})

	t.Run("Should override endpoints set by environment variables with custom endpoints in code", func(st *testing.T) {
		setUp()
		defer cleanup()
		c := &Configuration{}
		c.loadEnv()

		if got, exp := c.Endpoints.Notify, notifyEndpoint; got != exp {
			st.Errorf("Expected Notify endpoint: '%s', but was: '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, sessionsEndpoint; got != exp {
			st.Errorf("Expected Sessions endpoint: '%s', but was: '%s'", exp, got)
		}

		notifyOverride := "https://notify.override.com"
		sessionsOverride := "https://sessions.override.com"
		c.update(&Configuration{
			Endpoints: Endpoints{
				Notify:   notifyOverride,
				Sessions: sessionsOverride,
			},
		})
		if got, exp := c.Endpoints.Notify, notifyOverride; got != exp {
			st.Errorf("Expected Notify endpoint: '%s', but was: '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, sessionsOverride; got != exp {
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

func TestInsightHubEndpoints(t *testing.T) {
	hubNotify := "https://notify.insighthub.smartbear.com"
	hubSession := "https://sessions.insighthub.smartbear.com"
	customNofify := "https://custom.notify.com/"
	customSessions := "https://custom.sessions.com/"
	hubApiKey := "00000abcdef0123456789abcdef012345"

	setUp := func() (*Configuration, *CustomTestLogger) {
		logger := &CustomTestLogger{}
		return &Configuration{
			Logger: logger,
		}, logger
	}

	t.Run("Should use InsightHub endpoints if API key has prefix", func(st *testing.T) {
		c, _ := setUp()
		c.update(&Configuration{
			APIKey: hubApiKey,
		})

		if got, exp := c.Endpoints.Notify, hubNotify; got != exp {
			st.Errorf("Expected notify endpoint to be '%s' but was '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, hubSession; got != exp {
			st.Errorf("Expected sessions endpoint to be '%s' but was '%s'", exp, got)
		}
	})

	t.Run("Should prefer custom endpoints over InsightHub endpoints", func(st *testing.T) {
		c, _ := setUp()
		c.update(&Configuration{
			APIKey: hubApiKey,
			Endpoints: Endpoints{
				Notify:   customNofify,
				Sessions: customSessions,
			},
		})
		if got, exp := c.Endpoints.Notify, customNofify; got != exp {
			st.Errorf("Expected notify endpoint to be '%s' but was '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, customSessions; got != exp {
			st.Errorf("Expected sessions endpoint to be '%s' but was '%s'", exp, got)
		}
	})

	t.Run("With InsightHub API key and only custom notify endpoint, sessions should be empty", func(st *testing.T) {
		c, _ := setUp()
		c.update(&Configuration{
			APIKey: hubApiKey,
			Endpoints: Endpoints{
				Notify: customNofify,
			},
		})
		if got, exp := c.Endpoints.Notify, customNofify; got != exp {
			st.Errorf("Expected notify endpoint to be '%s' but was '%s'", exp, got)
		}
		if got, exp := c.Endpoints.Sessions, ""; got != exp {
			st.Errorf("Expected sessions endpoint to be empty but was '%s'", got)
		}
	})

	t.Run("With InsightHub API key and only custom session endpoint, panic should be thrown", func(st *testing.T) {
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
			APIKey: hubApiKey,
			Endpoints: Endpoints{
				Sessions: customSessions,
			},
		})
	})
}
