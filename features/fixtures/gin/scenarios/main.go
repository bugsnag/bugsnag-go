package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/gin-gonic/gin"

	"github.com/bugsnag/bugsnag-go/gin"
)

func main() {
	testcase := flag.String("case", "", "test case to run")
	flag.Parse()

	// Increase publish rate for testing
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 20

	switch *testcase {
	case "default":
		caseDefault()

	default:
		panic("No valid test case: " + *testcase)
	}
}

func newDefaultConfig() bugsnag.Configuration {
	return bugsnag.Configuration{
		APIKey: os.Getenv("API_KEY"),
		Endpoints: bugsnag.Endpoints{
			Notify:   os.Getenv("BUGSNAG_ENDPOINT"),
			Sessions: os.Getenv("BUGSNAG_ENDPOINT"),
		},
	}
}

func caseDefault() {
	g := gin.Default()
	g.Use(bugsnaggin.AutoNotify(newDefaultConfig()))

	g.GET("/basic", func(c *gin.Context) {
		bugsnag.Notify(fmt.Errorf("oops"))
	})
	g.GET("/", func(c *gin.Context) {
		log.Println("ping")
	})

	g.Run(":4511")
}
