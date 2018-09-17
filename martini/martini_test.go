package bugsnagmartini_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/martini"
	. "github.com/bugsnag/bugsnag-go/testutil"
	"github.com/go-martini/martini"
)

func performHandledError(notifier *bugsnag.Notifier) {
	ctx := bugsnag.StartSession(context.Background())
	notifier.Notify(ctx, fmt.Errorf("Ooopsie"), bugsnag.User{Id: "987zyx"})
}

func performUnhandledCrash() string {
	var a struct{}
	crash(a)
	return "ok"
}

func TestMartini(t *testing.T) {
	ts, reports := Setup()
	defer ts.Close()

	m := martini.Classic()

	userID := "1234abcd"
	m.Use(martini.Recovery())
	config := bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	}
	bugsnag.Configure(config)
	m.Use(bugsnagmartini.AutoNotify(bugsnag.User{Id: userID}))

	m.Get("/unhandled", performUnhandledCrash)
	m.Get("/handled", performHandledError)
	go m.RunOnAddr(":9077")

	t.Run("AutoNotify", func(st *testing.T) {
		time.Sleep(1 * time.Second)
		_, err := http.Get("http://localhost:9077/unhandled")
		if err != nil {
			t.Error(err)
		}
		report := <-reports
		r, _ := simplejson.NewJson(report)
		hostname, _ := os.Hostname()
		AssertPayload(st, r, fmt.Sprintf(`
		{
			"apiKey": "%s",
			"events":[
				{
					"app":{ "releaseStage":"" },
					"context":"/unhandled",
					"device":{ "hostname": "%s" },
					"exceptions":[
						{
							"errorClass":"*runtime.TypeAssertionError",
							"message":"interface conversion: interface {} is struct {}, not string",
							"stacktrace":[]
						}
					],
					"metaData":{
						"request":{ "httpMethod":"GET", "url":"http://localhost:9077/unhandled" }
					},
					"payloadVersion":"4",
					"severity":"error",
					"severityReason":{ "type":"unhandledErrorMiddleware" },
					"unhandled":true,
					"user":{ "id": "%s" }
				}
			],
			"notifier":{
				"name":"Bugsnag Go",
				"url":"https://github.com/bugsnag/bugsnag-go",
				"version": "%s"
			}
		}
		`, TestAPIKey, hostname, userID, bugsnag.VERSION))
	})

	t.Run("Notify", func(st *testing.T) {
		time.Sleep(1 * time.Second)
		_, err := http.Get("http://localhost:9077/handled")
		if err != nil {
			t.Error(err)
		}
		report := <-reports
		r, _ := simplejson.NewJson(report)
		hostname, _ := os.Hostname()
		AssertPayload(st, r, fmt.Sprintf(`
		{
			"apiKey": "%s",
			"events":[
				{
					"app":{ "releaseStage":"" },
					"device":{ "hostname": "%s" },
					"exceptions":[
						{
							"errorClass":"*errors.errorString",
							"message":"Ooopsie",
							"stacktrace":[]
						}
					],
					"metaData":{
						"request":{ "httpMethod":"GET", "url":"http://localhost:9077/handled" }
					},
					"payloadVersion":"4",
					"severity":"error",
					"severityReason":{ "type":"unhandledErrorMiddleware" },
					"unhandled":true,
					"user":{ "id": "%s" }
				}
			],
			"notifier":{
				"name":"Bugsnag Go",
				"url":"https://github.com/bugsnag/bugsnag-go",
				"version": "%s"
			}
		}
		`, TestAPIKey, hostname, "987zyx", bugsnag.VERSION))
	})

}

func crash(a interface{}) string {
	return a.(string)
}
