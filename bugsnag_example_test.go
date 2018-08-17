package bugsnag_test

import (
	"fmt"
	"net"
	"net/http"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

var testAPIKey = "166f5ad3590596f9aa8d601ea89af845"
var testEndpoint string

func ExampleAutoNotify() {
	createAccount := func() {
		fmt.Println("Creating account...")
	}

	//AutoNotify would report any panics that happen
	handlerFunc := func(w http.ResponseWriter, request *http.Request) {
		defer bugsnag.AutoNotify(request, bugsnag.Context{String: "createAccount"})
		createAccount()
	}

	var w http.ResponseWriter
	var request *http.Request
	handlerFunc(w, request)
	// Output:
	// Creating account...
}

func ExampleRecover() {
	panicFunc := func() {
		fmt.Println("About to panic")
		panic("Oh noes")
	}

	// Will recover when panicFunc panics
	func() {
		config := bugsnag.Configuration{APIKey: testAPIKey}
		defer bugsnag.Recover(config)
		panicFunc()
	}()

	fmt.Println("Panic recovered")
	// Output:
	// About to panic
	// Panic recovered
}

func ExampleConfigure() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       "YOUR_API_KEY_HERE",
		ReleaseStage: "production",
		// See Configuration{} for other fields
	})
}

func ExampleHandler() {
	handleGet := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling GET")
	}

	// Set up your http handlers as usual
	http.HandleFunc("/", handleGet)

	// use bugsnag.Handler(nil) to wrap the default http handlers
	// so that Bugsnag is automatically notified about panics.
	http.ListenAndServe(":1234", bugsnag.Handler(nil))
}

func ExampleHandler_customServer() {
	handleGet := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling GET")
	}

	// If you're using a custom server, set the handlers explicitly.
	http.HandleFunc("/", handleGet)

	srv := http.Server{
		Addr:        ":1234",
		ReadTimeout: 10 * time.Second,
		// use bugsnag.Handler(nil) to wrap the default http handlers
		// so that Bugsnag is automatically notified about panics.
		Handler: bugsnag.Handler(nil),
	}
	srv.ListenAndServe()
}

func ExampleHandler_customHandlers() {
	handleGet := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling GET")
	}

	// If you're using custom handlers, wrap the handlers explicitly.
	handler := http.NewServeMux()
	http.HandleFunc("/", handleGet)
	// use bugsnag.Handler(handler) to wrap the handlers so that Bugsnag is
	// automatically notified about panics
	http.ListenAndServe(":1234", bugsnag.Handler(handler))
}

func ExampleNotify() {
	_, err := net.Listen("tcp", ":80")

	if err != nil {
		bugsnag.Notify(err)
	}
}

func ExampleNotify_details() {
	_, err := net.Listen("tcp", ":80")

	if err != nil {
		bugsnag.Notify(err,
			// show as low-severity
			bugsnag.SeverityInfo,
			// set the context
			bugsnag.Context{String: "createlistener"},
			// pass the user id in to count users affected.
			bugsnag.User{Id: "123456789"},
			// custom meta-data tab
			bugsnag.MetaData{
				"Listen": {
					"Protocol": "tcp",
					"Port":     "80",
				},
			},
		)
	}
}

func ExampleOnBeforeNotify() {

	type Job struct {
		Retry     bool
		UserID    string
		UserEmail string
	}

	bugsnag.OnBeforeNotify(func(event *bugsnag.Event, config *bugsnag.Configuration) error {
		// Search all the RawData for any *Job pointers that we're passed in
		// to bugsnag.Notify() and friends.
		for _, datum := range event.RawData {
			if job, ok := datum.(*Job); ok {
				// don't notify bugsnag about errors in retries
				if job.Retry {
					return fmt.Errorf("bugsnag middleware: not notifying about job retry")
				}
				// add the job as a tab on Bugsnag.com
				event.MetaData.AddStruct("Job", job)
				// set the user correctly
				event.User = &bugsnag.User{Id: job.UserID, Email: job.UserEmail}
			}
		}

		// continue notifying as normal
		return nil
	})
}
