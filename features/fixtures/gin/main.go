package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/gin-gonic/gin"

	"github.com/bugsnag/bugsnag-go/gin"
)

func main() {
	g := gin.Default()
	config := bugsnag.Configuration{
		APIKey: os.Getenv("API_KEY"),
		Endpoints: bugsnag.Endpoints{
			Notify:   os.Getenv("BUGSNAG_ENDPOINT"),
			Sessions: os.Getenv("BUGSNAG_ENDPOINT"),
		},
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

	acs, err := strconv.ParseBool(os.Getenv("AUTO_CAPTURE_SESSIONS"))
	if err == nil {
		config.AutoCaptureSessions = acs
	}
	bugsnag.Configure(config)

	// Increase publish rate for testing
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 300

	g.Use(gin.Recovery(), bugsnaggin.AutoNotify(config))

	g.GET("/autonotify-then-recover", unhandledCrash)
	g.GET("/handled", handledError)
	g.GET("/session", session)
	g.GET("/autonotify", autonotify)
	g.GET("/onbeforenotify", onBeforeNotify)
	g.GET("/recover", dontDie)
	g.GET("/user", user)
	g.Run(":" + os.Getenv("SERVER_PORT"))

}

func unhandledCrash(c *gin.Context) {
	// Invalid type assertion, will panic
	func(a interface{}) string {
		return a.(string)
	}(struct{}{})
}

func handledError(c *gin.Context) {
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		if errClass := os.Getenv("ERROR_CLASS"); errClass != "" {
			bugsnag.Notify(err, c.Request.Context(), bugsnag.ErrorClass{Name: errClass})
		} else {
			bugsnag.Notify(err, c.Request.Context())
		}
	}
}

func session(c *gin.Context) {
	log.Println("single session")
}

func dontDie(c *gin.Context) {
	defer bugsnag.Recover(c.Request.Context())
	panic("Request killed but recovered")
}

func user(c *gin.Context) {
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
		Id:    "test-user-id",
		Name:  "test-user-name",
		Email: "test-user-email",
	})
}

func onBeforeNotify(c *gin.Context) {
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
	time.Sleep(100 * time.Millisecond)
	bugsnag.Notify(fmt.Errorf("Don't ignore this error"))
	time.Sleep(100 * time.Millisecond)
	bugsnag.Notify(fmt.Errorf("Change error message"))
	time.Sleep(100 * time.Millisecond)
}

func autonotify(c *gin.Context) {
	go func(ctx context.Context) {
		defer func() { recover() }()
		defer bugsnag.AutoNotify(ctx)
		panic("Go routine killed with auto notify")
	}(c.Request.Context())
}
