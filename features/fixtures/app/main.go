package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bugsnag/bugsnag-go/v2"
)

var scenariosMap = map[string] func()(bugsnag.Configuration, func()){
	"UnhandledScenario": UnhandledCrashScenario,
	"HandledScenario": HandledErrorScenario,
	"MultipleUnhandledScenario": MultipleUnhandledErrorsScenario,
	"MultipleHandledScenario": MultipleHandledErrorsScenario,
	"NestedErrorScenario": NestedHandledErrorScenario,
	"MetadataScenario": MetadataScenario,
	"FilteredMetadataScenario": FilteredMetadataScenario,
	"HandledCallbackErrorScenario": HandledCallbackErrorScenario,
	"SendSessionScenario": SendSessionScenario,
	"HandledToUnhandledScenario": HandledToUnhandledScenario,
	"SetUserScenario": SetUserScenario,
	"RecoverAfterPanicScenario": RecoverAfterPanicScenario,
	"AutonotifyPanicScenario": AutonotifyPanicScenario,
	"SessionAndErrorScenario": SessionAndErrorScenario,
	"OnBeforeNotifyScenario": OnBeforeNotifyScenario,
}

func main() {
		// Listening to the OS Signals
		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		ticker := time.NewTicker(1 * time.Second)

		addr := os.Getenv("DEFAULT_MAZE_ADDRESS")
		if (addr == "") {
			addr = DEFAULT_MAZE_ADDRESS
		}
	
		for {
			select {
			case <-ticker.C:
				fmt.Println("[Bugsnag] Get command")
				command := GetCommand(DEFAULT_MAZE_ADDRESS)
				fmt.Printf("[Bugsnag] Received command: %+v\n", command)

				if command.Action == "run-scenario" {
					prepareScenarioFunc, ok := scenariosMap[command.ScenarioName]
					if ok {
						config, scenarioFunc := prepareScenarioFunc()
						bugsnag.Configure(config)
						scenarioFunc()
					}
				}
			case <-ctx.Done():
					fmt.Println("[Bugsnag] Context is done, closing")
					ticker.Stop()
					return
			}
		}
}


