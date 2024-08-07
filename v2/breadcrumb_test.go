package bugsnag_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/bugsnag/bugsnag-go/v2/testutil"
)

func TestDefaultBreadcrumbValues(t *testing.T) {
	testServer, reports, notifier := setupServer(bugsnag.Configuration{})
	defer testServer.Close()
	notifier.LeaveBreadcrumb("test breadcrumb")
	notifier.Notify(fmt.Errorf("test error"))
	breadcrumbs := getBreadcrumbs(reports)

	if len(breadcrumbs) != 1 {
		t.Fatal("expected 1 breadcrumb")
	}
	if breadcrumbs[0].Name != "test breadcrumb" {
		t.Fatal("expected breadcrumb name")
	}
	if len(breadcrumbs[0].Timestamp) < 6 {
		t.Fatal("expected timestamp")
	}
	if len(breadcrumbs[0].MetaData) != 0 {
		t.Fatal("expected no metadata")
	}
	if breadcrumbs[0].Type != bugsnag.BreadcrumbTypeManual {
		t.Fatal("expected manual type")
	}
}

func TestCustomBreadcrumbValues(t *testing.T) {
	testServer, reports, notifier := setupServer(bugsnag.Configuration{})
	defer testServer.Close()
	notifier.LeaveBreadcrumb("test breadcrumb", bugsnag.BreadcrumbMetaData{"hello": "world"}, bugsnag.BreadcrumbTypeProcess)
	notifier.Notify(fmt.Errorf("test error"))
	breadcrumbs := getBreadcrumbs(reports)

	if len(breadcrumbs) != 1 {
		t.Fatal("expected 1 breadcrumb")
	}
	if breadcrumbs[0].Name != "test breadcrumb" {
		t.Fatal("expected breadcrumb name")
	}
	if len(breadcrumbs[0].Timestamp) < 6 {
		t.Fatal("expected timestamp")
	}
	if len(breadcrumbs[0].MetaData) != 1 || breadcrumbs[0].MetaData["hello"] != "world" {
		t.Fatal("expected correct metadata")
	}
	if breadcrumbs[0].Type != bugsnag.BreadcrumbTypeProcess {
		t.Fatal("expected process type")
	}
}

func TestDefaultMaxBreadcrumbs(t *testing.T) {
	testServer, reports, notifier := setupServer(bugsnag.Configuration{})
	defer testServer.Close()
	defaultMaximum := 25

	for i := 1; i <= defaultMaximum*2; i++ {
		notifier.LeaveBreadcrumb(fmt.Sprintf("breadcrumb%v", i))
	}

	notifier.Notify(fmt.Errorf("test error"))
	breadcrumbs := getBreadcrumbs(reports)

	if len(breadcrumbs) != defaultMaximum {
		t.Fatal("incorrect number of breadcrumbs")
	}
	for i := 0; i < defaultMaximum; i++ {
		if breadcrumbs[i].Name != fmt.Sprintf("breadcrumb%v", defaultMaximum*2-i) {
			t.Fatal("invalid breadcrumb at ", i)
		}
	}
}

func TestCustomMaxBreadcrumbs(t *testing.T) {
	customMaximum := 5
	testServer, reports, notifier := setupServer(bugsnag.Configuration{MaximumBreadcrumbs: customMaximum})
	defer testServer.Close()

	for i := 1; i <= customMaximum*2; i++ {
		notifier.LeaveBreadcrumb(fmt.Sprintf("breadcrumb%v", i))
	}

	notifier.Notify(fmt.Errorf("test error"))
	breadcrumbs := getBreadcrumbs(reports)

	if len(breadcrumbs) != customMaximum {
		t.Fatal("incorrect number of breadcrumbs")
	}
	for i := 0; i < customMaximum; i++ {
		if breadcrumbs[i].Name != fmt.Sprintf("breadcrumb%v", customMaximum*2-i) {
			t.Fatal("invalid breadcrumb at ", i)
		}
	}
}

func TestBreadcrumbCallbacksAreReversed(t *testing.T) {
	testServer, reports, notifier := setupServer(bugsnag.Configuration{})
	defer testServer.Close()

	callback1Called := false
	callback2Called := false
	notifier.BreadcrumbState.OnBreadcrumb(func(breadcrumb *bugsnag.Breadcrumb) bool {
		callback2Called = true
		if breadcrumb.Name != "breadcrumb" {
			t.Fatal("incorrect name")
		}
		if callback1Called == false {
			t.Fatal("callbacks should occur in reverse order")
		}
		return true
	})
	notifier.BreadcrumbState.OnBreadcrumb(func(breadcrumb *bugsnag.Breadcrumb) bool {
		callback1Called = true
		if breadcrumb.Name != "breadcrumb" {
			t.Fatal("incorrect name")
		}
		if callback2Called == true {
			t.Fatal("callbacks should occur in reverse order")
		}
		return true
	})

	notifier.LeaveBreadcrumb("breadcrumb")

	if !callback2Called {
		t.Fatal("breadcrumb callback not called")
	}

	notifier.Notify(fmt.Errorf("test error"))
	if len(getBreadcrumbs(reports)) != 1 {
		t.Fatal("expected one breadcrumb")
	}
}

func TestBreadcrumbCallbacksCanCancel(t *testing.T) {
	testServer, reports, notifier := setupServer(bugsnag.Configuration{})
	defer testServer.Close()

	callbackCalled := false
	notifier.BreadcrumbState.OnBreadcrumb(func(breadcrumb *bugsnag.Breadcrumb) bool {
		t.Fatal("Callback should be canceled")
		return true
	})
	notifier.BreadcrumbState.OnBreadcrumb(func(breadcrumb *bugsnag.Breadcrumb) bool {
		callbackCalled = true
		return false
	})

	notifier.LeaveBreadcrumb("breadcrumb")

	if !callbackCalled {
		t.Fatal("first breadcrumb callback not called")
	}

	notifier.Notify(fmt.Errorf("test error"))
	if len(getBreadcrumbs(reports)) != 0 {
		t.Fatal("breadcrumb not canceled")
	}
}

func TestSendNoBreadcrumbs(t *testing.T) {
	testServer, reports, notifier := setupServer(bugsnag.Configuration{})
	defer testServer.Close()
	notifier.Notify(fmt.Errorf("test error"))
	if len(getBreadcrumbs(reports)) != 0 {
		t.Fatal("expected no breadcrumbs")
	}
}

func TestSendOrderedBreadcrumbs(t *testing.T) {
	testServer, reports, notifier := setupServer(bugsnag.Configuration{})
	defer testServer.Close()
	notifier.LeaveBreadcrumb("breadcrumb1")
	notifier.LeaveBreadcrumb("breadcrumb2")
	notifier.Notify(fmt.Errorf("test error"))
	breadcrumbs := getBreadcrumbs(reports)
	if len(breadcrumbs) != 2 {
		t.Fatal("expected 2 breadcrumbs", breadcrumbs)
	}
	if breadcrumbs[0].Name != "breadcrumb2" || breadcrumbs[1].Name != "breadcrumb1" {
		t.Fatal("expected ordered breadcrumbs", breadcrumbs)
	}
}

func TestSendCleanMetadata(t *testing.T) {
	testServer, reports, notifier := setupServer(bugsnag.Configuration{})
	defer testServer.Close()
	type Recursive struct {
		Inner *Recursive
	}
	recursiveValue := Recursive{}
	recursiveValue.Inner = &recursiveValue
	notifier.LeaveBreadcrumb("breadcrumb2", bugsnag.BreadcrumbMetaData{"recursive": recursiveValue})
	notifier.Notify(fmt.Errorf("test error"))
	breadcrumbs := getBreadcrumbs(reports)
	if len(breadcrumbs) != 1 {
		t.Fatal("expected 1 breadcrumb", breadcrumbs)
	}
	if breadcrumbs[0].MetaData["recursive"].(map[string]interface{})["Inner"] != "[RECURSION]" {
		t.Fatal("remove recursive")
	}
}

func getBreadcrumbs(reports chan []byte) []bugsnag.Breadcrumb {
	event, _ := simplejson.NewJson(<-reports)
	fmt.Println(event)
	firstEventJson := testutil.GetIndex(event, "events", 0)
	breadcrumbsJson := testutil.Get(firstEventJson, "breadcrumbs")

	breadcrumbs := []bugsnag.Breadcrumb{}
	for index := 0; index < len(breadcrumbsJson.MustArray()); index++ {
		breadcrumbJson := breadcrumbsJson.GetIndex(index)
		fmt.Println(breadcrumbJson)
		breadcrumbs = append(breadcrumbs, bugsnag.Breadcrumb{
			Timestamp: breadcrumbJson.Get("timestamp").MustString(),
			Name:      breadcrumbJson.Get("name").MustString(),
			Type:      breadcrumbJson.Get("type").MustString(),
			MetaData:  breadcrumbJson.Get("metaData").MustMap(),
		})
	}
	return breadcrumbs
}

func setupServer(configuration bugsnag.Configuration) (*httptest.Server, chan []byte, *bugsnag.Notifier) {
	testServer, reports := testutil.Setup()
	configuration.APIKey = testutil.TestAPIKey
	configuration.Endpoints = bugsnag.Endpoints{Notify: testServer.URL, Sessions: testServer.URL + "/sessions"}
	notifier := bugsnag.New(configuration)
	return testServer, reports, notifier
}
