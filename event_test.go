package bugsnag

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestPopulateEvent(t *testing.T) {
	event := new(Event)
	contexts := make(chan context.Context, 1)
	reqs := make(chan *http.Request, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contexts <- AttachRequestData(r.Context(), r)
		reqs <- r
	}))
	defer ts.Close()

	http.Get(ts.URL + "/serenity?q=abcdef")

	ctx, req := <-contexts, <-reqs
	populateEventWithContext(ctx, event)

	for _, tc := range []struct{ e, c interface{} }{
		{e: event.Ctx, c: ctx},
		{e: event.Request, c: extractRequestInfoFromReq(req)},
		{e: event.Context, c: req.URL.Path},
		{e: event.User.Id, c: req.RemoteAddr[:strings.LastIndex(req.RemoteAddr, ":")]},
	} {
		if !reflect.DeepEqual(tc.e, tc.c) {
			t.Errorf("Expected '%+v' and '%+v' to be equal", tc.e, tc.c)
		}
	}
}
