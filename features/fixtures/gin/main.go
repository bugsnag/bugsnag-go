package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	config := bugsnag.Configuration{
		AppVersion: os.Getenv("APP_VERSION"),
		AppType:    os.Getenv("APP_TYPE"),
		APIKey:     os.Getenv("API_KEY"),
		Endpoint:   os.Getenv("ENDPOINT"),
		Endpoints: bugsnag.Endpoints{
			Notify:   os.Getenv("NOTIFY_ENDPOINT"),
			Sessions: os.Getenv("SESSIONS_ENDPOINT"),
		},
		Hostname:     os.Getenv("HOSTNAME"),
		ReleaseStage: os.Getenv("RELEASE_STAGE"),
	}

	if stages := os.Getenv("NOTIFY_RELEASE_STAGES"); stages != "" {
		config.NotifyReleaseStages = []string{stages}
	}

	if acs, _ := strconv.ParseBool(os.Getenv("AUTO_CAPTURE_SESSIONS")); acs {
		config.AutoCaptureSessions = acs
	}

	if filters := os.Getenv("PARAMS_FILTERS"); filters != "" {
		config.ParamsFilters = []string{filters}

	}

	config.Synchronous, _ = strconv.ParseBool(os.Getenv("SYNCHRONOUS"))
	bugsnag.Configure(config)

	g.Use(bugsnaggin.AutoNotify(config))

	g.GET("/unhandled", unhandledCrash)
	g.GET("/handled", handledError)
	g.GET("/metadata", metadata)
	g.GET("/onbeforenotify", onbeforenotify)
	g.GET("/recover", dontdie)
	g.GET("/async", async)
	g.GET("/user", user)

	g.Run(":9050") // listen and serve on 0.0.0.0:9001
}

func unhandledCrash(c *gin.Context) {
	// Invalid type assertion, will panic
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func handledError(c *gin.Context) {
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		bugsnag.Notify(err, c.Request.Context())
	}
}

func metadata(c *gin.Context) {
	customerData := map[string]string{"Name": "Joe Bloggs", "Age": "21"}
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
		"Scheme": {
			"Customer": customerData,
			"Level":    "Blue",
		},
	})
}

func dontdie(c *gin.Context) {
	defer bugsnag.Recover()
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func async(c *gin.Context) {
	bugsnag.Notify(fmt.Errorf("If I show up it means I was sent synchronously"))
	defer os.Exit(0)
}

func user(c *gin.Context) {
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
		Id:    "test-user-id",
		Name:  "test-user-name",
		Email: "test-user-email",
	})
}

func onbeforenotify(c *gin.Context) {
	bugsnag.OnBeforeNotify(
		func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			if event.Message == "Ignore this error" {
				return fmt.Errorf("not sending errors to ignore")
			}
			// continue notifying as normal
			if event.Message == "Change error message" {
				event.Message = "Error message was changed"
			}
			return nil
		})
	bugsnag.Notify(fmt.Errorf("Ignore this error"))
	bugsnag.Notify(fmt.Errorf("Don't ignore this error"))
	bugsnag.Notify(fmt.Errorf("Change error message"))
}
