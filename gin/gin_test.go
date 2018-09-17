package bugsnaggin_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/gin"
	. "github.com/bugsnag/bugsnag-go/testutil"
	"github.com/gin-gonic/gin"
)

const (
	testAPIKey = "166f5ad3590596f9aa8d601ea89af845"
	port       = "9079"
)

// setup sets up and returns a test event server for receiving the event payloads.
// report payloads published to the returned server's URL will be put on the returned channel
func setup() (*httptest.Server, chan []byte) {
	reports := make(chan []byte, 10)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Bugsnag called")
		if strings.Contains(r.URL.Path, "sessions") {
			return
		}
		body, _ := ioutil.ReadAll(r.Body)
		reports <- body
	})), reports
}

func TestGin(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	g := gin.Default()

	userID := "1234abcd"
	g.Use(bugsnaggin.AutoNotify(bugsnag.Configuration{
		APIKey:    testAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	}, bugsnag.User{Id: userID}))

	g.GET("/unhandled", performUnhandledCrash)
	g.GET("/handled", performHandledError)
	go g.Run("localhost:9079") //This call blocks

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
							"errorClass":"*runtime.TypeAssertionError",
							"message":"interface conversion: interface {} is struct {}, not string",
							"stacktrace":[]
						}
					],
					"metaData":{
						"request":{ "httpMethod":"GET", "url":"http://localhost:9079/unhandled" }
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
	ctx := bugsnag.StartSession(context.Background())
	bugsnag.Notify(ctx, fmt.Errorf("Ooopsie"), bugsnag.User{Id: "987zyx"})
}

func performUnhandledCrash(c *gin.Context) {
	c.String(http.StatusOK, "OK")
	var a struct{}
	crash(a)
}

func crash(a interface{}) string {
	return a.(string)
}
