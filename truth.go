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

func handleHeaders(w http.ResponseWriter, r *http.Request) {
	responseJSON, err := json.Marshal(collapseMapVals(r.Header))
	if err != nil {
		logerror(err.Error())
		return
	}
	io.WriteString(w, string(responseJSON)+"\n")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	response := map[string]bool{"running": true}
	responseJSON, _ := json.Marshal(response)
	io.WriteString(w, string(responseJSON)+"\n")
}

func main() {
	loadRoutes()
	server()
}

func loadRoutes() {
	routes = make(map[string]func(http.ResponseWriter, *http.Request))
	routes["/kubernetes/canary"] = handlePing
	routes["/headers"] = handleHeaders
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

	loginfo("starting up server on port " + port)
	server.ListenAndServe()
}

// http.Request.Header is a map where the values are arrays of strings,
// and it's nicer to look at output where there's a single key/value
func collapseMapVals(input map[string][]string) map[string]string {
	result := map[string]string{}
	for key, value := range input {
		result[key] = value[0]
	}
	return result
}

type requestLog struct {
	Timestamp   string
	HTTPCode    int
	RequestPath string
	RemoteHost  string
	Headers     map[string]string
}

func loginfo(msg string) {
	fmt.Fprintf(os.Stderr, "{\"info\":\"%s\"}\n", msg)
}

func logerror(msg string) {
	fmt.Fprintf(os.Stderr, "{\"error\":\"%s\"}\n", msg)
}

func log(r *http.Request, statusCode int) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	logline := requestLog{
		Timestamp:   time.Now().Format(time.RFC3339),
		HTTPCode:    statusCode,
		RequestPath: r.URL.Path,
		RemoteHost:  ip,
		Headers:     collapseMapVals(r.Header),
	}

	json, err := json.Marshal(logline)
	if err != nil {
		logerror(err.Error())
		return
	}

	fmt.Println(string(json))
}
