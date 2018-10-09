package bugsnagrevel

import (
	"strings"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/revel/revel"
)

type headerExposer interface {
	GetAll(string) []string
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
func extractHeaders(h headerExposer, paramsFilters []string) map[string]string {
	m := make(map[string]string)
	// Standard headers
	addHeader(paramsFilters, m, h, "Accept")
	addHeader(paramsFilters, m, h, "Accept-Charset")
	addHeader(paramsFilters, m, h, "Accept-Encoding")
	addHeader(paramsFilters, m, h, "Accept-Language")
	addHeader(paramsFilters, m, h, "Accept-Datetime")
	addHeader(paramsFilters, m, h, "Access-Control-Request-Method")
	addHeader(paramsFilters, m, h, "Access-Control-Request-Headers")
	addHeader(paramsFilters, m, h, "Authorization")
	addHeader(paramsFilters, m, h, "Cache-Control")
	addHeader(paramsFilters, m, h, "Connection")
	addHeader(paramsFilters, m, h, "Content-Length")
	addHeader(paramsFilters, m, h, "Content-Type")
	addHeader(paramsFilters, m, h, "Cookie")
	addHeader(paramsFilters, m, h, "Date")
	addHeader(paramsFilters, m, h, "Expect")
	addHeader(paramsFilters, m, h, "Forwarded")
	addHeader(paramsFilters, m, h, "From")
	addHeader(paramsFilters, m, h, "Host")
	addHeader(paramsFilters, m, h, "If-Match")
	addHeader(paramsFilters, m, h, "If-Modified-Since")
	addHeader(paramsFilters, m, h, "If-None-Match")
	addHeader(paramsFilters, m, h, "If-Range")
	addHeader(paramsFilters, m, h, "If-Unmodified-Since")
	addHeader(paramsFilters, m, h, "Max-Forwards")
	addHeader(paramsFilters, m, h, "Origin")
	addHeader(paramsFilters, m, h, "Pragma")
	addHeader(paramsFilters, m, h, "Proxy-Authorization")
	addHeader(paramsFilters, m, h, "Range")
	addHeader(paramsFilters, m, h, "Referer")
	addHeader(paramsFilters, m, h, "TE")
	addHeader(paramsFilters, m, h, "User-Agent")
	addHeader(paramsFilters, m, h, "Upgrade")
	addHeader(paramsFilters, m, h, "Via")
	addHeader(paramsFilters, m, h, "Warning")

	// Non-standard but common
	addHeader(paramsFilters, m, h, "DNT")
	addHeader(paramsFilters, m, h, "X-Requested-With")
	addHeader(paramsFilters, m, h, "X-CSRF-Token")
	return m
}

func addHeader(paramsFilters []string, m map[string]string, h headerExposer, headerName string) {
	if val := h.GetAll(headerName); val != nil {
		if contains(paramsFilters, strings.ToLower(headerName)) {
			m[headerName] = "[REDACTED]"
		} else {
			m[headerName] = strings.Join(val, ",")
		}
	}
}

func contains(slice []string, sub string) bool {
	for _, s := range slice {
		if s == sub {
			return true
		}
	}
	return false
}
