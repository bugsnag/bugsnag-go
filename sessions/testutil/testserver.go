package testutil

import (
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

var sessionTestOnce sync.Once

// SessionRequest embodies the payload and headers received in a request to the
// sessions endpoint
type SessionRequest struct {
	Body   []byte
	Header http.Header
}

// StartSessionTestServer starts up a mock session server for testing on the
// given authority. Returns a channel that will write any requests (body +
// header) it receives
func StartSessionTestServer(sessionAuthority string) chan SessionRequest {
	c := make(chan SessionRequest, 10)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		c <- SessionRequest{Body: body, Header: r.Header}
	})
	l, err := net.Listen("tcp", sessionAuthority)
	if err != nil {
		panic(err)
	}
	go http.Serve(l, mux)
	return c
}
