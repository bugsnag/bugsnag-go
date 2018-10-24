package main

import (
	"fmt"
	"os"

	"github.com/bugsnag/bugsnag-go"
)

func main() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey: os.Getenv("API_KEY"),
		Endpoints: bugsnag.Endpoints{
			Notify:   os.Getenv("NOTIFY_ENDPOINT"),
			Sessions: os.Getenv("SESSIONS_ENDPOINT"),
		},
	})
	notifier := bugsnag.New()
	notifier.NotifySync(fmt.Errorf("oops"), true)
}
