package bugsnag

import (
	"context"
	"net/http"
	"strings"
)

const (
	requestContextKey     requestKey = iota + 1
	requestJSONContextKey requestKey = iota + 1
)

type requestKey int

// AttachRequestData returns a child of the given context with the request
// object attached for later extraction by the notifier in order to
// automatically record request data
func AttachRequestData(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestContextKey, r)
}

// AttachRequestJSONData returns a child of the given context with the request
// JSON object attached for later extraction by the notifier in order to
// automatically record request data. Similar to AttachRequestData, but expects
// that the request data is already extracted into a *bugsnag.RequestJSON
// instead of a later extracting from a *http.Request
func AttachRequestJSONData(ctx context.Context, json *RequestJSON) context.Context {
	return context.WithValue(ctx, requestJSONContextKey, json)
}

// extractRequestInfo looks for the request object that the notifier
// automatically attaches to the context when using any of the supported
// frameworks or bugsnag.HandlerFunc or bugsnag.Handler, and returns sub-object
// supported by the notify API.
func extractRequestInfo(ctx context.Context) *RequestJSON {
	if req := getRequestJSONIfPresent(ctx); req != nil {
		return req
	}
	if req := getRequestIfPresent(ctx); req != nil {
		return extractRequestInfoFromReq(req)
	}
	return nil
}

// extractRequestInfoFromReq extracts the request information the notify API
// understands from the given HTTP request. Returns the sub-object supported by
// the notify API.
func extractRequestInfoFromReq(req *http.Request) *RequestJSON {
	return &RequestJSON{
		ClientIP:   req.RemoteAddr,
		HTTPMethod: req.Method,
		URL:        req.RequestURI,
		Referer:    req.Referer(),
		Headers:    parseRequestHeaders(req.Header),
	}
}

func parseRequestHeaders(header map[string][]string) map[string]string {
	headers := make(map[string]string)
	for k, v := range header {
		// Headers can have multiple values, in which case we report them as csv
		if contains(Config.ParamsFilters, k) {
			headers[k] = "[FILTERED]"
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

// Certain frameworks (Revel) don't use the standard library to model HTTP
// requests. In this case we extract the request information upfront and place
// this sub-payload in the context directly. If for some reason the context
// contains both a *http.Request object and a *bugsnag.RequestJSON then the
// *bugsnag.RequestJSON will take priority.
func getRequestJSONIfPresent(ctx context.Context) *RequestJSON {
	if ctx == nil {
		return nil
	}
	val := ctx.Value(requestJSONContextKey)
	if val == nil {
		return nil
	}
	return val.(*RequestJSON)
}
