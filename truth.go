// Minimal server to respond to requests with JSON
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var routes map[string]func(http.ResponseWriter, *http.Request)

var port = "4242"

const httpReturnRequestedCodePath = "/httpCode/"

type router struct{}

func handleHTTPCode(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.URL.Path[len(httpReturnRequestedCodePath):])
	log(r, n)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if n == 301 {
		w.Header().Set("Location", "/foobar")
	}
	w.WriteHeader(n)
	io.WriteString(w, http.StatusText(n))
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	response := map[string]bool{"running": true}
	responseJSON, _ := json.Marshal(response)
	io.WriteString(w, string(responseJSON))
}

func main() {
	loadRoutes()
	server()
}

func loadRoutes() {
	routes = make(map[string]func(http.ResponseWriter, *http.Request))
	routes["/kubernetes/canary"] = handlePing
	routes["/ping"] = handlePing
	routes["/"] = handlePing
}

func (*router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler, ok := routes[request.URL.String()]
	if ok {
		log(request, 200)
		handler(writer, request)
	} else if strings.HasPrefix(request.URL.String(), httpReturnRequestedCodePath) {
		// handle "/httpCode/XXX" - return requested httpCode to the client
		handleHTTPCode(writer, request)
	} else {
		log(request, 404)
		http.Error(writer, "404 page not found", 404)
	}
}

func server() {
	server := http.Server{
		Addr:    ":" + port,
		Handler: &router{},
	}

	fmt.Println("starting up server on port " + port)
	server.ListenAndServe()
}

func log(r *http.Request, statusCode int) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	forwarded := r.Header.Get("X-Forwarded-For")
	fmt.Fprintf(
		os.Stdout,
		"{\"timestamp\":\"%s\",\"http_code\":%d,\"request_path\":\"%s\",\"remote_host\":\"%s\",\"x-forwarded_for\":\"%s\"}\n",
		time.Now().Format(time.RFC3339),
		statusCode,
		r.URL.Path,
		ip,
		forwarded,
	)
}
