package bugsnagrevel

import (
	"testing"

	"github.com/bugsnag/bugsnag-go"
)

type testExposer map[string]string

func (te testExposer) GetAll(n string) []string {
	return []string{te[n]}
}

func TestExtractingHeaders(t *testing.T) {
	var te testExposer = make(map[string]string)
	te["Accept"] = "application/json"
	te["Accept-Charset"] = "utf-8"
	te["Accept-Encoding"] = "gzip"
	te["Accept-Language"] = "en-US"
	te["Accept-Datetime"] = "Thu, 31 May 2007 20:35:00 GMT"
	te["Access-Control-Request-Method"] = "GET"
	te["Cache-Control"] = "no-cache"
	te["Connection"] = "keep-alive"
	te["Content-Length"] = "42"
	te["Content-Type"] = "application/x-www-form-urlencoded"
	te["Date"] = "Tue, 15 Nov 1994 08:12:31 GMT"
	te["Expect"] = "100-continue"
	te["Forwarded"] = "for=192.0.2.60; proto=http; by=203.0.113.43"
	te["From"] = "user@example.com"
	te["Host"] = "bugsnag.com"
	te["If-Match"] = "737060cd8c284d8582d"
	te["If-Modified-Since"] = "Sat, 29 Oct 1994 19:43:31 GMT"
	te["If-None-Match"] = "737060cd8c284d8582d"
	te["If-Range"] = "737060cd8c284d8582d"
	te["If-Unmodified-Since"] = "Sat, 29 Oct 1994 19:43:31 GMT"
	te["Max-Forwards"] = "29"
	te["Origin"] = "https://bugsnag.com"
	te["Pragma"] = "no-cache"
	te["Proxy-Authorization"] = "Basic 2323jiojioIJOIOJIJ=="
	te["Range"] = "bytes=400-999"
	te["Referer"] = "https://bugsnag.com"
	te["TE"] = "trailers"
	te["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"
	te["Upgrade"] = "h2c, HTTPS/1.3, IRC/6.9, RTA/x11, websocket"
	te["Via"] = "1.0 fred, 1.1 bugsnag.com (Apache/1.1)"
	te["Warning"] = "199 Miscellaneous warning"
	te["DNT"] = "1"
	te["X-Requested-With"] = "XMLHttpRequest"
	te["X-CSRF-Token"] = "<TOKEN>"
	m := extractHeaders(te, bugsnag.Config.ParamsFilters)
	for k, v := range te {
		if got, exp := m[k], v; exp != got {
			t.Errorf("Expected '%s' to be '%s' but was '%s'", k, exp, got)
		}
	}

	key := "Authorization"
	te[key] = "Basic 34i3j4iom2323=="
	m = extractHeaders(te, bugsnag.Config.ParamsFilters)
	if got, exp := m[key], "[REDACTED]"; got != exp {
		t.Errorf("Expected '%s' to be '%s' but was '%s'", key, exp, got)
	}

	key = "Cookie"
	te[key] = "name=value"
	m = extractHeaders(te, bugsnag.Config.ParamsFilters)
	if got, exp := m[key], "[REDACTED]"; got != exp {
		t.Errorf("Expected '%s' to be '%s' but was '%s'", key, exp, got)
	}
}

type testMultiValExposer map[string][]string

func (tmve testMultiValExposer) GetAll(n string) []string {
	return tmve[n]
}

func TestSupportsMultipleValues(t *testing.T) {
	var tmve testMultiValExposer = make(map[string][]string)
	key := "Access-Control-Request-Headers"
	tmve[key] = []string{"origin", "x-requested-with", "accept"}
	m := extractHeaders(tmve, []string{})
	if got, exp := m[key], "origin,x-requested-with,accept"; got != exp {
		t.Errorf("Expected '%s' to be '%s' but was '%s'", key, exp, got)
	}
}
