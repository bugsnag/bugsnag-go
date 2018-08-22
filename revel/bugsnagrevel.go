// Package bugsnagrevel adds Bugsnag to revel.
// It lets you pass *revel.Controller into bugsnag.Notify(),
// and provides a Filter to catch errors.
package bugsnagrevel

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/bugsnag/bugsnag-go"
	"github.com/revel/revel"
)

var once sync.Once

// FrameworkName is the name of the framework this middleware applies to
const FrameworkName string = "Revel"

var errorHandlingState = bugsnag.HandledState{
	SeverityReason:   bugsnag.SeverityReasonUnhandledMiddlewareError,
	OriginalSeverity: bugsnag.SeverityError,
	Unhandled:        true,
	Framework:        FrameworkName,
}

// Filter should be added to the filter chain just after the PanicFilter.
// It sends errors to Bugsnag automatically. Configuration is read out of
// conf/app.conf, you should set bugsnag.apikey, and can also set
// bugsnag.endpoints, bugsnag.releasestage, bugsnag.apptype, bugsnag.appversion,
// bugsnag.projectroot, bugsnag.projectpackages if needed.
func Filter(c *revel.Controller, fc []revel.Filter) {
	// Record a session if auto capture sessions is enabled
	if bugsnag.Config.IsAutoCaptureSessions() {
		ctx := bugsnag.StartSession(context.Background())
		c.Args["context"] = ctx
	}

	defer bugsnag.AutoNotify(c, errorHandlingState)
	fc[0](c, fc[1:])
}

// Add support to bugsnag for reading data out of *revel.Controllers
func middleware(event *bugsnag.Event, config *bugsnag.Configuration) error {
	for _, datum := range event.RawData {
		if controller, ok := datum.(*revel.Controller); ok {
			// make the request visible to the builtin HttpMiddleware
			if version("0.18.0") {
				event.RawData = append(event.RawData, controller.Request)
			} else {
				req := struct{ *http.Request }{}
				event.RawData = append(event.RawData, req.Request)
			}
			event.Context = controller.Action
			event.MetaData.AddStruct("Session", controller.Session)
		}
	}

	return nil
}

func init() {
	revel.OnAppStart(func() {
		bugsnag.OnBeforeNotify(middleware)

		var projectPackages []string
		if packages, ok := revel.Config.String("bugsnag.projectpackages"); ok {
			projectPackages = strings.Split(packages, ",")
		} else {
			projectPackages = []string{revel.ImportPath + "/app/*", revel.ImportPath + "/app"}
		}

		bugsnag.Configure(bugsnag.Configuration{
			APIKey:              revel.Config.StringDefault("bugsnag.apikey", ""),
			AutoCaptureSessions: revel.Config.BoolDefault("bugnsnag.autocapturesessions", true),
			Endpoints: bugsnag.Endpoints{
				Notify:   revel.Config.StringDefault("bugsnag.endpoints.notify", ""),
				Sessions: revel.Config.StringDefault("bugsnag.endpoints.sessions", ""),
			},
			AppType:         revel.Config.StringDefault("bugsnag.apptype", ""),
			AppVersion:      revel.Config.StringDefault("bugsnag.appversion", ""),
			ReleaseStage:    revel.Config.StringDefault("bugsnag.releasestage", revel.RunMode),
			ProjectPackages: projectPackages,
			Logger:          revel.ERROR,
		})
	})
}

// Very basic semantic versioning.
// Returns true if given version matches or is above revel.Version
func version(reqVersion string) bool {
	req := strings.Split(reqVersion, ".")
	cur := strings.Split(revel.Version, ".")
	for i := 0; i < 2; i++ {
		rV, _ := strconv.Atoi(req[i])
		cV, _ := strconv.Atoi(cur[i])
		if (rV < cV && i == 0) || (rV < cV && i == 1) {
			return true
		}
	}
	return false
}
