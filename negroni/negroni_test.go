package bugsnagnegroni_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/negroni"
	. "github.com/bugsnag/bugsnag-go/testutil"
	"github.com/urfave/negroni"
)

const userID = "1234abcd"

func TestNegroni(t *testing.T) {
	ts, reports := Setup()
	config := bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/unhandled", unhandledCrashHandler)
	mux.HandleFunc("/handled", handledCrashHandler)

	hostname, _ := os.Hostname()

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(bugsnagnegroni.AutoNotify(config, bugsnag.User{Id: userID}))
	n.UseHandler(mux)

	go http.ListenAndServe(":9078", n)

	t.Run("AutoNotify", func(st *testing.T) {
		time.Sleep(500 * time.Millisecond)
		http.Get("http://localhost:9078/unhandled")
		report := <-reports
		r, _ := simplejson.NewJson(report)
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
							"message":"something went terribly wrong",
							"stacktrace":[]
						}
					],
					"payloadVersion":"4",
					"severity":"error",
					"severityReason":{ "type":"unhandledErrorMiddleware" },
					"unhandled":true,
					"request": {
						"url": "http://localhost:9078/unhandled",
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

	t.Run("Notify", func(st *testing.T) {
		time.Sleep(500 * time.Millisecond)
		http.Get("http://localhost:9078/handled")
		report := <-reports
		r, _ := simplejson.NewJson(report)
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
					"unhandled": false,
					"request": {
						"url": "http://localhost:9078/handled",
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
}

func unhandledCrashHandler(w http.ResponseWriter, req *http.Request) {
	panic("something went terribly wrong")
}

func handledCrashHandler(w http.ResponseWriter, req *http.Request) {
	bugsnag.Notify(fmt.Errorf("Ooopsie"), bugsnag.User{Id: userID}, req.Context())
}

func crash(a interface{}) string {
	return a.(string)
}
