package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bugsnag/bugsnag-go/v2"
)

var scenariosMap = map[string]func(Command) func(){
	"UnhandledScenario":            UnhandledCrashScenario,
	"HandledScenario":              HandledErrorScenario,
	"MultipleUnhandledScenario":    MultipleUnhandledErrorsScenario,
	"MultipleHandledScenario":      MultipleHandledErrorsScenario,
	"NestedErrorScenario":          NestedHandledErrorScenario,
	"MetadataScenario":             MetadataScenario,
	"FilteredMetadataScenario":     FilteredMetadataScenario,
	"HandledCallbackErrorScenario": HandledCallbackErrorScenario,
	"SendSessionScenario":          SendSessionScenario,
	"HandledToUnhandledScenario":   HandledToUnhandledScenario,
	"SetUserScenario":              SetUserScenario,
	"RecoverAfterPanicScenario":    RecoverAfterPanicScenario,
	"AutonotifyPanicScenario":      AutonotifyPanicScenario,
	"SessionAndErrorScenario":      SessionAndErrorScenario,
	"OnBeforeNotifyScenario":       OnBeforeNotifyScenario,
	"AutoconfigPanicScenario":      AutoconfigPanicScenario,
	"AutoconfigHandledScenario":    AutoconfigHandledScenario,
	"AutoconfigMetadataScenario":   AutoconfigMetadataScenario,
	"HttpServerScenario":           HttpServerScenario,
}

func main() {
	addr := os.Getenv("DEFAULT_MAZE_ADDRESS")
	if addr == "" {
		addr = DEFAULT_MAZE_ADDRESS
	}

	endpoints := bugsnag.Endpoints{
		Notify:   fmt.Sprintf("%+v/notify", addr),
		Sessions: fmt.Sprintf("%+v/sessions", addr),
	}
	// HAS TO RUN FIRST BECAUSE OF PANIC WRAP
	// https://github.com/bugsnag/panicwrap/blob/master/panicwrap.go#L177-L203
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:    "166f5ad3590596f9aa8d601ea89af845",
		Endpoints: endpoints,
	})
	// Increase publish rate for testing
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 50

	// Listening to the OS Signals
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			command := GetCommand(addr)
			fmt.Printf("[Bugsnag] Received command: %+v\n", command)
			if command.Action != "run-scenario" {
				continue
			}
			prepareScenarioFunc, ok := scenariosMap[command.ScenarioName]
			if ok {
				scenarioFunc := prepareScenarioFunc(command)
				scenarioFunc()
				time.Sleep(200 * time.Millisecond)
			}
		}
	}
}
