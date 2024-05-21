package bugsnag_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bugsnag/bugsnag-go/v2"
	. "github.com/bugsnag/bugsnag-go/v2/testutil"
)

const NUM_RUNS = 100

func benchmarkServerSetup() (*httptest.Server) {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

func benchmarkNotifierSetup(url string) *bugsnag.Notifier {
	return bugsnag.New(bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: url, Sessions: url + "/sessions"},
	})
}

func BenchmarkStandardAsyncNotify(b *testing.B) {
	ts := benchmarkServerSetup()
	defer ts.Close()
	notifier := benchmarkNotifierSetup(ts.URL)

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	})

	for i:=0; i<b.N; i++ {
		for j:=0; j<NUM_RUNS; j++ {
			msg := fmt.Sprintf("Oopsie %+v", fmt.Sprint(j))
			notifier.Notify(fmt.Errorf(msg))
		}
	}

	time.Sleep(10 *time.Second)
}

func BenchmarkStandardSyncNotify(b *testing.B) {
	ts := benchmarkServerSetup()
	defer ts.Close()
	notifier := benchmarkNotifierSetup(ts.URL)

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	})

	for i:=0; i<b.N; i++ {
		for j:=0; j<NUM_RUNS; j++ {
			msg := fmt.Sprintf("Oopsie %+v", fmt.Sprint(j))
			notifier.NotifySync(fmt.Errorf(msg), true)
		}
	}

	time.Sleep(10 *time.Second)
}

func BenchmarkStandardAsyncPoolNotify(b *testing.B) {
	ts := benchmarkServerSetup()
	defer ts.Close()
	notifier := benchmarkNotifierSetup(ts.URL)

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	})

	for i:=0; i<b.N; i++ {
		for j:=0; j<NUM_RUNS; j++ {
			msg := fmt.Sprintf("Oopsie %+v", fmt.Sprint(j))
			notifier.NotifyAsyncPool(fmt.Errorf(msg), true)
		}
	}

	time.Sleep(10 *time.Second)
}
