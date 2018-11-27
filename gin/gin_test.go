package bugsnaggin_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/gin"
	. "github.com/bugsnag/bugsnag-go/testutil"
	"github.com/gin-gonic/gin"
)

func TestGin(t *testing.T) {
	ts, reports := Setup()
	defer ts.Close()

	g := gin.Default()

	userID := "1234abcd"
	g.Use(bugsnaggin.AutoNotify(bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	}, bugsnag.User{Id: userID}))

	g.GET("/unhandled", performUnhandledCrash)
	g.GET("/handled", performHandledError)
	go g.Run(":9079") //This call blocks

	t.Run("AutoNotify", func(st *testing.T) {
		time.Sleep(1 * time.Second)
		_, err := http.Get("http://localhost:9079/unhandled")
		if err != nil {
			t.Error(err)
		}
		report := <-reports
		r, _ := simplejson.NewJson(report)
		hostname, _ := os.Hostname()
		AssertPayload(st, r, fmt.Sprintf(`
		{
			"apiKey":"166f5ad3590596f9aa8d601ea89af845",
			"events":[
				{
					"app":{ "releaseStage":"" },
					"context":"/unhandled",
					"device":{ "hostname": "%s" },
					"exceptions":[
						{
							"errorClass":"*errors.errorString",
							"message":"you shouldn't have done that",
							"stacktrace":[]
						}
					],
					"payloadVersion":"4",
					"severity":"error",
					"severityReason":{ "type":"unhandledErrorMiddleware" },
					"unhandled":true,
					"request": {
						"url": "http://localhost:9079/unhandled",
						"httpMethod": "GET",
						"referer": "",
						"headers": {
							"Accept-Encoding": "gzip"
						}
					},
					"user":{ "id": "%s" }
				}
			],
			"notifier":{
				"name":"Bugsnag Go",
				"url":"https://github.com/bugsnag/bugsnag-go",
				"version": "%s"
			}
		}
		`, hostname, userID, bugsnag.VERSION))
	})

	t.Run("Manual notify", func(st *testing.T) {
		_, err := http.Get("http://localhost:9079/handled")
		if err != nil {
			t.Error(err)
		}
		report := <-reports
		r, _ := simplejson.NewJson(report)
		hostname, _ := os.Hostname()
		AssertPayload(st, r, fmt.Sprintf(`
		{
			"apiKey":"166f5ad3590596f9aa8d601ea89af845",
			"events":[
				{
					"app":{ "releaseStage":"" },
					"context":"/handled",
					"device":{ "hostname": "%s" },
					"exceptions":[
						{
							"errorClass":"*errors.errorString",
							"message":"Ooopsie",
							"stacktrace":[]
						}
					],
					"payloadVersion":"4",
					"severity":"warning",
					"severityReason":{ "type":"handledError" },
					"unhandled":false,
					"request": {
						"url": "http://localhost:9079/handled",
						"httpMethod": "GET",
						"referer": "",
						"headers": {
							"Accept-Encoding": "gzip"
						}
					},
					"user":{ "id": "%s" }
				}
			],
			"notifier":{
				"name":"Bugsnag Go",
				"url":"https://github.com/bugsnag/bugsnag-go",
				"version": "%s"
			}
		}
		`, hostname, "987zyx", bugsnag.VERSION))
	})
}

func performHandledError(c *gin.Context) {
	ctx := c.Request.Context()
	bugsnag.Notify(fmt.Errorf("Ooopsie"), ctx, bugsnag.User{Id: "987zyx"})
}

func performUnhandledCrash(c *gin.Context) {
	panic("you shouldn't have done that")
}

func crash(a interface{}) string {
	return a.(string)
}
