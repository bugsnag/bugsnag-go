package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bugsnag/bugsnag-go/v2"
)

// Insert your API key
const apiKey = "YOUR-API-KEY-HERE"

func main() {
	if len(apiKey) != 32 {
		fmt.Println("Please set the API key in main.go before running the example")
		return
	}

	bugsnag.Configure(bugsnag.Configuration{APIKey: apiKey})

	http.HandleFunc("/unhandled", unhandledCrash)
	http.HandleFunc("/handled", handledError)

	fmt.Println("=============================================================================")
	fmt.Println("Visit http://localhost:9001/unhandled - To perform an unhandled crash")
	fmt.Println("Visit http://localhost:9001/handled   - To create a manual error notification")
	fmt.Println("=============================================================================")
	fmt.Println("")

	http.ListenAndServe(":9001", bugsnag.Handler(nil))
}

func unhandledCrash(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))

	// Invalid type assertion, will panic
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func handledError(w http.ResponseWriter, r *http.Request) {
	_, err := os.Open("nonexistent_file.txt")
	if err != nil {
		bugsnag.Notify(err, r.Context())
	}
}
