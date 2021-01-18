package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/bugsnag/bugsnag-go/v2"
)

// Insert your API key
const apiKey = "YOUR-API-KEY-HERE"

// The following example will cause two events in your dashboard:
// One event because AutoNotify intercepted a panic.
// The other because Bugsnag noticed your application was about to be taken
// down by a panic.
// To avoid taking down your application and the last event, replace
// bugsnag.AutoNotify with bugsnag.Recover in the below example.
func main() {
	if len(apiKey) != 32 {
		fmt.Println("Please set your API key in main.go before running example.")
		return
	}

	bugsnag.Configure(bugsnag.Configuration{APIKey: apiKey})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		fmt.Println("Starting new go routine...")
		// Manually create a new Bugsnag session for this goroutine
		ctx := bugsnag.StartSession(context.Background())
		defer wg.Done()
		// AutoNotify captures any panics, repanicking after error reports are sent
		defer bugsnag.AutoNotify(ctx)

		// Invalid type assertion, will panic
		func(a interface{}) { _ = a.(string) }(struct{}{})
	}()

	wg.Wait()
}
