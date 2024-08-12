package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bugsnag/bugsnag-go/v2"
)

var scenariosMap = map[string] func(Command)(bugsnag.Configuration, func()){
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
	"AutoconfigPanicScenario": AutoconfigPanicScenario,
	"AutoconfigHandledScenario": AutoconfigHandledScenario,
	"AutoconfigMetadataScenario": AutoconfigMetadataScenario,
}

func main() {
		// Listening to the OS Signals
		signalsChan := make(chan os.Signal, 1)
		signal.Notify(signalsChan, syscall.SIGINT, syscall.SIGTERM)
		ticker := time.NewTicker(1 * time.Second)

		// Increase publish rate for testing
		bugsnag.DefaultSessionPublishInterval = time.Millisecond * 50

		addr := os.Getenv("DEFAULT_MAZE_ADDRESS")
		if (addr == "") {
			addr = DEFAULT_MAZE_ADDRESS
		}
	
		for {
			select {
			case <-ticker.C:
				command := GetCommand(addr)
				fmt.Printf("[Bugsnag] Received command: %+v\n", command)

				if command.Action == "run-scenario" {
					prepareScenarioFunc, ok := scenariosMap[command.ScenarioName]
					if ok {
						config, scenarioFunc := prepareScenarioFunc(command)
						bugsnag.Configure(config)
						time.Sleep(200 * time.Millisecond)
						scenarioFunc()
						time.Sleep(200 * time.Millisecond)
					}
				}
			case <-signalsChan:
					fmt.Println("[Bugsnag] Context is done, closing")
					ticker.Stop()
					return
			}
		}
}


