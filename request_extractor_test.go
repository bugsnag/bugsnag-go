package bugsnag

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestRequestInformationGetsExtracted(t *testing.T) {
	contexts := make(chan context.Context, 1)
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = AttachRequestData(ctx, r)
		contexts <- ctx
	})
	ts := httptest.NewServer(hf)
	defer ts.Close()
	http.Get(ts.URL + "/1234abcd")

	reqJSON := extractRequestInfo(<-contexts)
	if reqJSON.ClientIP == "" {
		t.Errorf("expected to find an IP address for the request but was blank")
	}
	if got, exp := reqJSON.HTTPMethod, "GET"; got != exp {
		t.Errorf("expected HTTP method to be '%s' but was '%s'", exp, got)
	}
	if got, exp := reqJSON.URL, "/1234abcd"; got != exp {
		t.Errorf("expected request URL to be '%s' but was '%s'", exp, got)
	}
	if got, exp := reqJSON.Referer, ""; got != exp {
		t.Errorf("expected request referer to be '%s' but was '%s'", exp, got)
	}
	if got, exp := reqJSON.Headers["Accept-Encoding"], "gzip"; got != exp {
		t.Errorf("expected Accept-Encoding to be '%s' but was '%s'", exp, got)
	}
	if got, exp := reqJSON.Headers["User-Agent"], "Go-http-client"; !strings.Contains(got, exp) {
		t.Errorf("expected user agent to contain '%s' but was '%s'", exp, got)
	}
}

func TestExtractingRequestJSON(t *testing.T) {
	json := &RequestJSON{
		ClientIP: "8.8.8.8",
		Headers: map[string]string{
			"most-goals": "Peter Crouch",
		},
		HTTPMethod: "GET",
		URL:        "/my-name-is-url",
		Referer:    "twitter",
	}
	ctx := AttachRequestJSONData(context.Background(), json)
	reqJSON := extractRequestInfo(ctx)
	if !reflect.DeepEqual(reqJSON, json) {
		t.Errorf("Expected JSON object '%+v' to be identical to '%+v'", reqJSON, json)
	}
}

func TestRequestExtractorCanHandleAbsentContext(t *testing.T) {
	if got := extractRequestInfo(nil); got != nil {
		//really just testing that nothing panics here
		t.Errorf("expected nil contexts to give nil sub-objects, but was '%s'", got)
	}
	if got := extractRequestInfo(context.Background()); got != nil {
		//really just testing that nothing panics here
		t.Errorf("expected contexts without requst info to give nil sub-objects, but was '%s'", got)
	}
}

func TestParseHeadersWillSanitiseIllegalParams(t *testing.T) {
	headers := make(map[string][]string)
	headers["password"] = []string{"correct horse battery staple"}
	headers["secret"] = []string{"I am Banksy"}
	headers["authorization"] = []string{"licence to kill -9"}
	for k, v := range parseRequestHeaders(headers) {
		if v != "[REDACTED]" {
			t.Errorf("expected '%s' to be [REDACTED], but was '%s'", k, v)
		}
	}
}
