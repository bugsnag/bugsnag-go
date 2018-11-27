package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/martini"
	"github.com/go-martini/martini"
)

// Insert your API key
const apiKey = "YOUR-API-KEY-HERE"

func main() {
	if len(apiKey) != 32 {
		fmt.Println("Please set the API key in main.go before running the example")
		return
	}

	bugsnag.Configure(bugsnag.Configuration{APIKey: apiKey})
	m := martini.Classic()

	m.Use(martini.Recovery())
	// Add bugsnag handler after martini.Recovery() to ensure panics get picked up
	m.Use(bugsnagmartini.AutoNotify())

	m.Get("/unhandled", performUnhandledCrash)
	m.Get("/handled", performHandledError)

	fmt.Println("=============================================================================")
	fmt.Println("Visit http://localhost:9001/unhandled - To perform an unhandled crash")
	fmt.Println("Visit http://localhost:9001/handled   - To create a manual error notification")
	fmt.Println("=============================================================================")
	fmt.Println("")

	m.RunOnAddr(":9001")
}

func performUnhandledCrash() {
	// Invalid type assertion, will panic
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func performHandledError(r *http.Request) {
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		bugsnag.Notify(err, r.Context())
	}
}
