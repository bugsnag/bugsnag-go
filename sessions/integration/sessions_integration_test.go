package sessions_test

import (
	"net"
	"net/http"
	"sync"
)

var authority = "localhost:8273"
var receivedRequests chan *http.Request
var sessionTestOnce sync.Once

var c = make(chan *http.Request, 10)

// StartServer starts up a mock session server for testing on the
// given authority. Returns a channel that will write any requests (body +
// header) it receives
func StartServer(sessionAuthority string) chan *http.Request {
	sessionTestOnce.Do(func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			c <- r
		})
		l, err := net.Listen("tcp", sessionAuthority)
		if err != nil {
			panic(err)
		}
		go http.Serve(l, mux)
	})
	return c
}
