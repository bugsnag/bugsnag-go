package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/negroni"
	"github.com/urfave/negroni"
)

// Insert your API key
const apiKey = "YOUR API KEY"

func main() {
	if len(apiKey) != 32 {
		fmt.Println("Please set your API key in main.go before running example.")
		return
	}

	bugsnag.Configure(bugsnag.Configuration{APIKey: apiKey})

	mux := http.NewServeMux()
	mux.HandleFunc("/unhandled", unhandledCrash)
	mux.HandleFunc("/handled", handledError)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	// Add bugsnag handler after negroni.NewRecovery() to ensure panics get picked up
	n.Use(bugsnagnegroni.AutoNotify())
	n.UseHandler(mux)

	fmt.Println("=============================================================================")
	fmt.Println("Visit http://localhost:9001/unhandled - To perform an unhandled crash")
	fmt.Println("Visit http://localhost:9001/handled   - To create a manual error notification")
	fmt.Println("=============================================================================")
	fmt.Println("")

	http.ListenAndServe(":9001", n)
}

func unhandledCrash(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))

	// Invalid type assertion, will panic
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func handledError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))
	_, err := os.Open("nonexistent_file.txt")
	if err != nil {
		bugsnag.Notify(err, r.Context())
	}
}
