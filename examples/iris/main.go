package main

import (
	"fmt"
	"os"

	"github.com/bugsnag/bugsnag-go"
	bugsnagiris "github.com/bugsnag/bugsnag-go/iris"

	"github.com/kataras/iris/v12"
)

// Insert your API key
const apiKey = "YOUR-API-KEY-HERE"

func main() {
	if len(apiKey) != 32 {
		fmt.Println("Please set your API key in main.go before running example.")
		return
	}

	app := iris.Default()

	app.Use(bugsnagiris.AutoNotify(bugsnag.Configuration{APIKey: apiKey}))

	app.Get("/unhandled", performUnhandledCrash)
	app.Get("/handled", performHandledError)

	fmt.Println("=============================================================================")
	fmt.Println("Visit http://localhost:9001/unhandled - To perform an unhandled crash")
	fmt.Println("Visit http://localhost:9001/handled   - To create a manual error notification")
	fmt.Println("=============================================================================")
	fmt.Println("")

	app.Listen(":9001") // listen and serve on 0.0.0.0:9001
}

func performUnhandledCrash(c iris.Context) {
	c.WriteString("OK")
	// Invalid type assertion, will panic
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func performHandledError(c iris.Context) {
	c.WriteString("OK")
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		bugsnag.Notify(err, c.Request.Context())
	}
}
