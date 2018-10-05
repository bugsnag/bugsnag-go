package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bugsnag/bugsnag-go"
)

func main() {
	http.HandleFunc("/unhandled", unhandledCrash)
	http.HandleFunc("/handledError", handledError)

	// Insert your API key
	bugsnag.Configure(bugsnag.Configuration{
		APIKey: "YOUR-API-KEY-HERE",
	})

	log.Println("Serving on 9001")
	http.ListenAndServe(":9001", bugsnag.Handler(nil))
}

func unhandledCrash(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))

	var a struct{}
	crash(a)
}

func handledError(w http.ResponseWriter, r *http.Request) {
	_, err := os.Open("some_nonexistent_file.txt")
	if err != nil {
		bugsnag.Notify(r.Context(), err)
	}
}

func crash(a interface{}) string {
	return a.(string)
}
