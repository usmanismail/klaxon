package router

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"
	"techtraits.com/log"
)

//A set of routes configured for this application
var routes []Route

// Registers a new handler
// Will check that it does not conflict with a route that is already configured
func Register(path string, method Method, consumes []string, produces []string,
	handler func(Request) (int, []byte)) bool {

	//Trim trailing slash if it exits
	var trailingSlash, _ = regexp.MatchString(".+/$", path)
	if trailingSlash {
		path = path[0 : len(path)-1]
	}

	regRoute := Route{path, method, consumes, produces, handler}
	for _, route := range routes {
		var match, status, _, _ = route.matchRoute(regRoute)
		if match && status == http.StatusOK {
			log.Error("Route %v already registered, ignoring ", regRoute)
			return false
		}
	}

	log.Debug("Route %v registered.", regRoute)
	routes = append(routes, regRoute)
	return true
}

func init() {
	http.HandleFunc("/", handler)
}

//Handles all incomming requests and routes them to registered handlers
func handler(resp http.ResponseWriter, req *http.Request) {

	reqRoute := Route{req.URL.Path, ToMethod(req.Method), strings.Split(strings.Split(req.Header.Get("Content-Type"), ";")[0], ","),
		strings.Split(strings.Split(req.Header.Get("Accept"), ";")[0], ","), nil}

	var httpStatusCode = http.StatusNotFound
	var errorMessage string = ""
	for _, route := range routes {
		var match, status, pathParams, message = route.matchRoute(reqRoute)
		if match && status == http.StatusOK {
			req.ParseForm()
			statusCode, content := route.Handler(RequestStruct{route, pathParams, req, resp})
			handleResponse(resp, statusCode, content)
			return
		} else if match && status >= httpStatusCode {
			httpStatusCode = status
			errorMessage = message
		}
	}
	handleResponse(resp, httpStatusCode, []byte(errorMessage))
}

func handleResponse(resp http.ResponseWriter, httpStatusCode int, message []byte) {

	if message == nil {
		message = []byte(http.StatusText(httpStatusCode))
	}

	if httpStatusCode == http.StatusOK {
		io.Copy(resp, bytes.NewReader(message))
	} else {
		http.Error(resp, string(message), httpStatusCode)
	}

}
