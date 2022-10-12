package bugsnag

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/bugsnag/bugsnag-go/errors"
	"github.com/bugsnag/bugsnag-go/sessions"
)

const expSmall = `{"apiKey":"","events":[{"app":{"releaseStage":""},"device":{"osName":"%s","runtimeVersions":{"go":"%s"}},"exceptions":[{"errorClass":"","message":"","stacktrace":null}],"metaData":{},"payloadVersion":"4","severity":"","unhandled":false}],"notifier":{"name":"Bugsnag Go","url":"https://github.com/bugsnag/bugsnag-go","version":"1.9.1"}}`

// The large payload has a timestamp in it which makes it awkward to assert against.
// Instead, assert that the timestamp property exist, along with the rest of the expected payload
const expLargePre = `{"apiKey":"166f5ad3590596f9aa8d601ea89af845","events":[{"app":{"releaseStage":"mega-production","type":"gin","version":"1.5.3"},"context":"/api/v2/albums","device":{"hostname":"super.duper.site","osName":"%s","runtimeVersions":{"go":"%s"}},"exceptions":[{"errorClass":"error class","message":"error message goes here","stacktrace":[{"method":"doA","file":"a.go","lineNumber":65},{"method":"fetchB","file":"b.go","lineNumber":99,"inProject":true},{"method":"incrementI","file":"i.go","lineNumber":651}]}],"groupingHash":"custom grouping hash","metaData":{"custom tab":{"my key":"my value"}},"payloadVersion":"4","session":{"startedAt":"`
const expLargePost = `,"severity":"info","severityReason":{"type":"unhandledError","attributes":{"framework":"gin"}},"unhandled":true,"user":{"id":"1234baerg134","name":"Kool Kidz on da bus","email":"typo@busgang.com"}}],"notifier":{"name":"Bugsnag Go","url":"https://github.com/bugsnag/bugsnag-go","version":"1.9.1"}}`

func TestMarshalEmptyPayload(t *testing.T) {
	sessionTracker = sessions.NewSessionTracker(&sessionTrackingConfig)
	p := payload{&Event{Ctx: context.Background()}, &Configuration{}}
	bytes, _ := p.MarshalJSON()
	exp := fmt.Sprintf(expSmall, runtime.GOOS, runtime.Version())
	if got := string(bytes[:]); got != exp {
		t.Errorf("Payload different to what was expected. \nGot: %s\nExp: %s", got, exp)
	}
}

func TestMarshalLargePayload(t *testing.T) {
	payload := makeLargePayload()
	bytes, _ := payload.MarshalJSON()
	got := string(bytes[:])
	expPre := fmt.Sprintf(expLargePre, runtime.GOOS, runtime.Version())
	if !strings.Contains(got, expPre) {
		t.Errorf("Expected large payload to contain\n'%s'\n but was\n'%s'", expPre, got)

	}
	if !strings.Contains(got, expLargePost) {
		t.Errorf("Expected large payload to contain\n'%s'\n but was\n'%s'", expLargePost, got)
	}
}

func makeLargePayload() *payload {
	stackframes := []StackFrame{
		{Method: "doA", File: "a.go", LineNumber: 65, InProject: false},
		{Method: "fetchB", File: "b.go", LineNumber: 99, InProject: true},
		{Method: "incrementI", File: "i.go", LineNumber: 651, InProject: false},
	}
	user := User{
		Id:    "1234baerg134",
		Name:  "Kool Kidz on da bus",
		Email: "typo@busgang.com",
	}
	handledState := HandledState{
		SeverityReason:   SeverityReasonUnhandledError,
		OriginalSeverity: severity{String: "error"},
		Unhandled:        true,
		Framework:        "gin",
	}

	ctx := context.Background()
	ctx = StartSession(ctx)

	event := Event{
		Error:        &errors.Error{},
		RawData:      nil,
		ErrorClass:   "error class",
		Message:      "error message goes here",
		Stacktrace:   stackframes,
		Context:      "/api/v2/albums",
		Severity:     SeverityInfo,
		GroupingHash: "custom grouping hash",
		User:         &user,
		Ctx:          ctx,
		MetaData: map[string]map[string]interface{}{
			"custom tab": map[string]interface{}{
				"my key": "my value",
			},
		},
		Unhandled: true,
		handledState: handledState,
	}
	config := Configuration{
		APIKey:       testAPIKey,
		ReleaseStage: "mega-production",
		AppType:      "gin",
		AppVersion:   "1.5.3",
		Hostname:     "super.duper.site",
	}
	return &payload{&event, &config}
}
