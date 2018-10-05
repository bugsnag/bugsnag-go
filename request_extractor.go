package bugsnag

import (
	"context"
	"net/http"
	"strings"
)

const requestContextKey requestKey = 0

type requestKey int

// extractRequestInfo looks for the request object that the notifier
// automatically attaches to the context when using any of the supported
// frameworks or bugsnag.HandlerFunc or bugsnag.Handler, and returns sub-object
// supported by the notify API.
func extractRequestInfo(ctx context.Context) *requestJSON {
	if req := getRequestIfPresent(ctx); req != nil {
		return extractRequestInfoFromReq(req)
	}
	return nil
}

// extractRequestInfoFromReq extracts the request information the notify API
// understands from the given HTTP request. Returns the sub-object supported by
// the notify API.
func extractRequestInfoFromReq(req *http.Request) *requestJSON {
	return &requestJSON{
		ClientIP:   req.RemoteAddr,
		HTTPMethod: req.Method,
		URL:        req.RequestURI,
		Referer:    req.Referer(),
		Headers:    parseRequestHeaders(req.Header),
	}
}

// attachRequestData returns a child of the given context with the request
// object attached for later extraction by the notifier in order to
// automatically record request data
func attachRequestData(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestContextKey, r)
}

func parseRequestHeaders(header map[string][]string) map[string]string {
	headers := make(map[string]string)
	for k, v := range header {
		// Headers can have multiple values, in which case we report them as csv
		if contains(Config.ParamsFilters, k) {
			headers[k] = "[REDACTED]"
		} else {
			headers[k] = strings.Join(v, ",")
		}
	}
	return headers
}

func contains(slice []string, e string) bool {
	for _, s := range slice {
		if s == e {
			return true
		}
	}
	return false
}
func getRequestIfPresent(ctx context.Context) *http.Request {
	if ctx == nil {
		return nil
	}
	val := ctx.Value(requestContextKey)
	if val == nil {
		return nil
	}
	return val.(*http.Request)
}
