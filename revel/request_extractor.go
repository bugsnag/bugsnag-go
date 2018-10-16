package bugsnagrevel

import (
	"strings"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/revel/revel"
)

type headerExposer interface {
	GetAll(string) []string
}

var headers = []string{
	// Standard headers
	"Accept",
	"Accept-Charset",
	"Accept-Encoding",
	"Accept-Language",
	"Accept-Datetime",
	"Access-Control-Request-Method",
	"Access-Control-Request-Headers",
	"Authorization",
	"Cache-Control",
	"Connection",
	"Content-Length",
	"Content-Type",
	"Cookie",
	"Date",
	"Expect",
	"Forwarded",
	"From",
	"Host",
	"If-Match",
	"If-Modified-Since",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
	"Max-Forwards",
	"Origin",
	"Pragma",
	"Proxy-Authorization",
	"Range",
	"Referer",
	"TE",
	"User-Agent",
	"Upgrade",
	"Via",
	"Warning",
	// Non-standard but common
	"DNT",
	"X-Requested-With",
	"X-CSRF-Token",
}

func extractRequestData(req *revel.Request, paramsFilters []string) *bugsnag.RequestJSON {
	return &bugsnag.RequestJSON{
		ClientIP: req.RemoteAddr,
		// Revel requires you to know the name of the headers up front, which
		// makes this bit tricky, especially if the application uses custom
		// headers
		Headers:    extractHeaders(req.Header, paramsFilters),
		HTTPMethod: req.Method,
		URL:        req.URL.String(),
		Referer:    req.Header.Get("referer"),
	}
}

// extractHeader is a best-guess workaround for pulling out as many headers as possible in a Revel request.
// Extracts all the standard HTTP headers, and some non-standard, but common headers.
func extractHeaders(exposer headerExposer, paramsFilters []string) map[string]string {
	m := make(map[string]string)
	for _, header := range headers {
		if val := exposer.GetAll(header); val != nil {
			if contains(paramsFilters, strings.ToLower(header)) {
				m[header] = "[FILTERED]"
			} else {
				m[header] = strings.Join(val, ",")
			}
		}
	}
	return m
}

func contains(slice []string, sub string) bool {
	for _, s := range slice {
		if s == sub {
			return true
		}
	}
	return false
}
