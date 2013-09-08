package router

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"techtraits.com/log"
)

//A set of routes configured for this application
var routes []Route

// Registers a new handler
// Will check that it does not conflict with a route that is already configured
func Register(path string, method Method, consumes []string, produces []string,
	handler func(Route, map[string]string, url.Values, http.Header)) bool {

	//Trim trailing slash if it exits
	var trailingSlash, _ = regexp.MatchString(".*/$", path)
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
			route.Handler(route, pathParams, req.Form, req.Header)
			return
		} else if match {
			httpStatusCode = status
			errorMessage = message
		}
	}
	handleError(resp, httpStatusCode, errorMessage)
}

func handleError(resp http.ResponseWriter, httpStatusCode int, message string) {

	if message == "" {
		switch httpStatusCode {
		case http.StatusBadRequest:
			message = "Bad Request"
		case http.StatusUnauthorized:
			message = "Unauthorized"
		case http.StatusPaymentRequired:
			message = "Payment Required"
		case http.StatusForbidden:
			message = "Forbidden"
		case http.StatusNotFound:
			message = "Resource not found"
		case http.StatusMethodNotAllowed:
			message = "Method not allowed"
		case http.StatusNotAcceptable:
			message = "Request not accepted"
		case http.StatusProxyAuthRequired:
			message = "Proxy auth required"
		case http.StatusRequestTimeout:
			message = "Request time out"
		case http.StatusConflict:
			message = "Conflict"
		case http.StatusGone:
			message = "Gone"
		case http.StatusLengthRequired:
			message = "Length Required"
		case http.StatusPreconditionFailed:
			message = "Preconditions failed"
		case http.StatusRequestEntityTooLarge:
			message = "Entity too large"
		case http.StatusRequestURITooLong:
			message = "URI too long"
		case http.StatusUnsupportedMediaType:
			message = "Unsupported media type"
		case http.StatusRequestedRangeNotSatisfiable:
			message = "Requested range not satisfiable"
		case http.StatusExpectationFailed:
			message = "Expectation failed"
		case http.StatusTeapot:
			message = "Unused"
		case http.StatusInternalServerError:
			message = "Internal Server Error"
		case http.StatusNotImplemented:
			message = "Not implemented"
		case http.StatusBadGateway:
			message = "Bad gateway"
		case http.StatusServiceUnavailable:
			message = "Service unavialable"
		case http.StatusGatewayTimeout:
			message = "Gateway timeout"
		case http.StatusHTTPVersionNotSupported:
			message = "Http version not supported"
		}
	}

	http.Error(resp, message, httpStatusCode)
}
