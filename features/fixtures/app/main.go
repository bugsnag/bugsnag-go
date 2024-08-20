package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
)

var scenariosMap = map[string]func(Command) func(){}

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

func multipleHandled() {
	//Make the order of the below predictable
	bugsnag.Configure(bugsnag.Configuration{Synchronous: true})

	ctx := bugsnag.StartSession(context.Background())
	bugsnag.Notify(fmt.Errorf("oops"), ctx)
	bugsnag.Notify(fmt.Errorf("oops"), ctx)
}

func multipleUnhandled() {
	//Make the order of the below predictable
	notifier := bugsnag.New(bugsnag.Configuration{Synchronous: true})
	notifier.FlushSessionsOnRepanic(false)

	ctx := bugsnag.StartSession(context.Background())
	defer func() { recover() }()
	defer notifier.AutoNotify(ctx)
	defer notifier.AutoNotify(ctx)
	panic("oops")
}

//go:noinline
func unhandledCrash() {
	// Invalid type assertion, will panic
	func(a interface{}) string {
		return a.(string)
	}(struct{}{})
}

func handledError() {
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		if errClass := os.Getenv("ERROR_CLASS"); errClass != "" {
			bugsnag.Notify(err, bugsnag.ErrorClass{Name: errClass})
		} else {
			bugsnag.Notify(err)
		}
	}
	// Give some time for the error to be sent before exiting
	time.Sleep(200 * time.Millisecond)
}

func session() {
	bugsnag.StartSession(context.Background())

	// Give some time for the session to be sent before exiting
	time.Sleep(200 * time.Millisecond)
}

func autonotify() {
	go func() {
		defer bugsnag.AutoNotify()
		panic("Go routine killed with auto notify")
	}()

	// Give enough time for the panic to happen
	time.Sleep(100 * time.Millisecond)
}

func metadata() {
	customerData := map[string]string{"Name": "Joe Bloggs", "Age": "21"}
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
		"Scheme": {
			"Customer": customerData,
			"Level":    "Blue",
		},
	})
	time.Sleep(200 * time.Millisecond)
}

func filtered() {
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
		"Account": {
			"Name":           "Company XYZ",
			"Price(dollars)": "1 Million",
		},
	})
	time.Sleep(200 * time.Millisecond)
}

func onBeforeNotify() {
	bugsnag.OnBeforeNotify(
		func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			if event.Message == "Ignore this error" {
				return fmt.Errorf("not sending errors to ignore")
			}
			// continue notifying as normal
			if event.Message == "Change error message" {
				event.Message = "Error message was changed"
			}
			return nil
		})
	bugsnag.Notify(fmt.Errorf("Ignore this error"))
	time.Sleep(100 * time.Millisecond)
	bugsnag.Notify(fmt.Errorf("Don't ignore this error"))
	time.Sleep(100 * time.Millisecond)
	bugsnag.Notify(fmt.Errorf("Change error message"))
	time.Sleep(100 * time.Millisecond)
}

func dontDie() {
	go func() {
		defer bugsnag.Recover()
		panic("Go routine killed but recovered")
	}()
	time.Sleep(100 * time.Millisecond)
}

func sessionAndError() {
	ctx := bugsnag.StartSession(context.Background())
	bugsnag.Notify(fmt.Errorf("oops"), ctx)

	time.Sleep(200 * time.Millisecond)
}

func sendAndExit() {
	bugsnag.Notify(fmt.Errorf("oops"))
}

func user() {
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
		Id:    "test-user-id",
		Name:  "test-user-name",
		Email: "test-user-email",
	})

	time.Sleep(200 * time.Millisecond)
}

func handledCallbackError() {
	bugsnag.Notify(fmt.Errorf("Inadequent Prep Error"), func(event *bugsnag.Event) {
		event.Context = "nonfatal.go:14"
		event.Severity = bugsnag.SeverityInfo

		event.Stacktrace[1].File = ">insertion<"
		event.Stacktrace[1].LineNumber = 0
	})
	// Give some time for the error to be sent before exiting
	time.Sleep(200 * time.Millisecond)
}

func handledToUnhandled() {
	bugsnag.Notify(fmt.Errorf("unknown event"), func(event *bugsnag.Event) {
		event.Unhandled = true
		event.Severity = bugsnag.SeverityError
	})
	// Give some time for the error to be sent before exiting
	time.Sleep(200 * time.Millisecond)
}

type customErr struct {
	msg     string
	cause   error
	callers []uintptr
}

func newCustomErr(msg string, cause error) error {
	callers := make([]uintptr, 8)
	runtime.Callers(2, callers)
	return customErr{
		msg:     msg,
		cause:   cause,
		callers: callers,
	}
}

func (err customErr) Error() string {
	return err.msg
}

func (err customErr) Unwrap() error {
	return err.cause
}

func (err customErr) Callers() []uintptr {
	return err.callers
}

func nestedHandledError() {
	if err := login("token " + os.Getenv("API_KEY")); err != nil {
		bugsnag.Notify(newCustomErr("terminate process", err))
		// Give some time for the error to be sent before exiting
		time.Sleep(200 * time.Millisecond)
	} else {
		i := len(os.Getenv("API_KEY"))
		// Some nonsense to avoid inlining checkValue
		if val, err := checkValue(i); err != nil {
			fmt.Printf("err: %v, val: %d", err, val)
		}
		if val, err := checkValue(i - 46); err != nil {
			fmt.Printf("err: %v, val: %d", err, val)
		}

		log.Fatalf("This test is broken - no error was generated.")
	}
}

func login(token string) error {
	val, err := checkValue(len(token) * -1)
	if err != nil {
		return newCustomErr("login failed", err)
	}
	fmt.Printf("val: %d", val)
	return nil
}

func checkValue(i int) (int, error) {
	if i < 0 {
		return 0, newCustomErr("invalid token", nil)
	} else if i%2 == 0 {
		return i / 2, nil
	} else if i < 9 {
		return i * 3, nil
	}

	return i * 4, nil
}
