package bugsnag

import (
	"testing"

	"github.com/bugsnag/bugsnag-go/errors"
)

const expSmall = `{"apiKey":"","events":[{"app":{"releaseStage":""},"device":{},"exceptions":[{"errorClass":"","message":"","stacktrace":null}],"metaData":{},"payloadVersion":"2","severity":"","unhandled":false}],"notifier":{"name":"Bugsnag Go","url":"https://github.com/bugsnag/bugsnag-go","version":"1.3.1"}}`
const expLarge = `{"apiKey":"166f5ad3590596f9aa8d601ea89af845","events":[{"app":{"releaseStage":"mega-production","type":"gin","version":"1.5.2"},"context":"/api/v2/albums","device":{"hostname":"super.duper.site"},"exceptions":[{"errorClass":"error class","message":"error message goes here","stacktrace":[{"method":"doA","file":"a.go","lineNumber":65},{"method":"fetchB","file":"b.go","lineNumber":99,"inProject":true},{"method":"incrementI","file":"i.go","lineNumber":651}]}],"groupingHash":"custom grouping hash","metaData":{"custom tab":{"my key":"my value"}},"payloadVersion":"2","severity":"info","severityReason":{"attributes":{"framework":"gin"},"type":"unhandledError"},"unhandled":true,"user":{"id":"1234baerg134","name":"Kool Kidz on da bus","email":"typo@busgang.com"}}],"notifier":{"name":"Bugsnag Go","url":"https://github.com/bugsnag/bugsnag-go","version":"1.3.1"}}`

func TestMarshalEmptyPayload(t *testing.T) {
	payload := payload{&Event{}, &Configuration{}}
	bytes, _ := payload.MarshalJSON()
	if got := string(bytes[:]); got != expSmall {
		t.Errorf("Payload different to what was expected. \nGot: %s\nExp: %s", got, expSmall)
	}
}

func TestMarshalLargePayload(t *testing.T) {
	payload := makeLargePayload()
	bytes, _ := payload.MarshalJSON()
	if got := string(bytes[:]); got != expLarge {
		t.Errorf("Payload different to what was expected. \nGot: %s\nExp: %s", got, expLarge)
	}
}

func BenchmarkEmptyMarshalJSON(b *testing.B) {
	payload := payload{&Event{}, &Configuration{}}
	for i := 0; i < b.N; i++ {
		payload.MarshalJSON()
	}
}

func BenchmarkLargeMarshalJSON(b *testing.B) {
	payload := makeLargePayload()
	for i := 0; i < b.N; i++ {
		payload.MarshalJSON()
	}
}

func makeLargePayload() *payload {
	stackframes := []stackFrame{
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
		MetaData: map[string]map[string]interface{}{
			"custom tab": map[string]interface{}{
				"my key": "my value",
			},
		},
		handledState: handledState}
	config := Configuration{
		APIKey:       testAPIKey,
		ReleaseStage: "mega-production",
		AppType:      "gin",
		AppVersion:   "1.5.2",
		Hostname:     "super.duper.site",
	}
	return &payload{&event, &config}
}
