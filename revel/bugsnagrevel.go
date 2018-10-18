// Package bugsnagrevel adds Bugsnag to revel.
// It lets you pass *revel.Controller into bugsnag.Notify(),
// and provides a Filter to catch errors.
package bugsnagrevel

import (
	"context"
	"net/http"
	"strings"

	"github.com/bugsnag/bugsnag-go"
	"github.com/revel/revel"
)

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
	notifier := bugsnag.New()
	ctx := bugsnag.AttachRequestData(context.Background(), findProperHTTPRequest(c))
	// Record a session if auto capture sessions is enabled
	if notifier.Config.IsAutoCaptureSessions() {
		ctx = bugsnag.StartSession(ctx)
	}
	c.Args["context"] = ctx
	defer notifier.AutoNotify(c, ctx, errorHandlingState)
	fc[0](c, fc[1:])
}

// Add support to bugsnag for reading data out of *revel.Controllers
func middleware(event *bugsnag.Event, config *bugsnag.Configuration) error {
	for _, datum := range event.RawData {
		if controller, ok := datum.(*revel.Controller); ok {
			// make the request visible to the builtin HttpMiddleware
			event.Context = controller.Action
			event.MetaData.AddStruct("Session", controller.Session)
		}
	}

	return nil
}

func findProperHTTPRequest(c *revel.Controller) *http.Request {
	var req *http.Request
	rawReq := c.Request.In.GetRaw()

	// This *should* always be a *http.Request, but the revel team must have
	// made this an interface{} for a reason, and we might as well be defensive
	// about it
	switch rawReq.(type) {
	case (*http.Request):
		req = rawReq.(*http.Request) // Find the *proper* http request.
	}
	return req
}

type bugsnagRevelLogger struct{}

func (l *bugsnagRevelLogger) Printf(s string, params ...interface{}) {
	if strings.HasPrefix(s, "ERROR") {
		revel.AppLog.Errorf(s, params...)
	} else if strings.HasPrefix(s, "WARN") {
		revel.AppLog.Warnf(s, params...)
	} else {
		revel.AppLog.Infof(s, params...)
	}

}

func init() {
	// To ensure that users can disable the default panic handler (by calling
	// bugsnag.Configure before this function does) we must allow other
	// callbacks to execute before this function.
	order := 2
	revel.OnAppStart(func() {
		bugsnag.OnBeforeNotify(middleware)

		ip := revel.ImportPath
		c := revel.Config
		bugsnag.Configure(bugsnag.Configuration{
			APIKey:   c.StringDefault("bugsnag.apikey", ""),
			Endpoint: c.StringDefault("bugsnag.endpoint", ""),
			Endpoints: bugsnag.Endpoints{
				Notify:   c.StringDefault("bugsnag.endpoints.notify", ""),
				Sessions: c.StringDefault("bugsnag.endpoints.sessions", ""),
			},
			ReleaseStage:        c.StringDefault("bugsnag.releasestage", revel.RunMode),
			AppType:             c.StringDefault("bugsnag.apptype", FrameworkName),
			AppVersion:          c.StringDefault("bugsnag.appversion", ""),
			AutoCaptureSessions: c.BoolDefault("bugsnag.autocapturesessions", true),
			Hostname:            c.StringDefault("bugsnag.device.hostname", ""),
			NotifyReleaseStages: getCsvsOrDefault("bugsnag.notifyreleasestages", nil),
			ProjectPackages:     getCsvsOrDefault("bugsnag.projectpackages", []string{ip + "/app/*", ip + "/app"}),
			SourceRoot:          c.StringDefault("bugsnag.sourceroot", ""),
			ParamsFilters:       getCsvsOrDefault("bugsnag.paramsfilters", []string{"password", "secret", "authorization", "cookie"}),
			Logger:              new(bugsnagRevelLogger),
			Synchronous:         c.BoolDefault("bugsnag.synchronous", false),
		})
	}, order)
}

func getCsvsOrDefault(propertyKey string, d []string) []string {
	if propString, ok := revel.Config.String(propertyKey); ok {
		return strings.Split(propString, ",")
	}
	return d
}
