package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/gin"
	"github.com/gin-gonic/gin"
)

// Insert your API key
const apiKey = "YOUR-API-KEY-HERE"

func main() {
	if len(apiKey) != 32 {
		fmt.Println("Please set your API key in main.go before running example.")
		return
	}

	g := gin.Default()

	g.Use(bugsnaggin.AutoNotify(bugsnag.Configuration{APIKey: apiKey}))

	g.GET("/unhandled", performUnhandledCrash)
	g.GET("/handled", performHandledError)

	fmt.Println("=============================================================================")
	fmt.Println("Visit http://localhost:9001/unhandled - To perform an unhandled crash")
	fmt.Println("Visit http://localhost:9001/handled   - To create a manual error notification")
	fmt.Println("=============================================================================")
	fmt.Println("")

	g.Run(":9001") // listen and serve on 0.0.0.0:9001
}

func performUnhandledCrash(c *gin.Context) {
	c.String(http.StatusOK, "OK")
	// Invalid type assertion, will panic
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func performHandledError(c *gin.Context) {
	c.String(http.StatusOK, "OK")
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		bugsnag.Notify(err, c.Request.Context())
	}
}
