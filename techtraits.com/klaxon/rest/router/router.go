package router

import (
   	"net/http"
   	"techtraits.com/log"
)

//A set of routes configured for this application
var routes []Route

// Registers a new handler
// Will check that it does not conflict with a route that is already configured
func Register(path string, method Method, consumes []string, produces []string, 
		handler func(Route, map[string]string, map[string]string, map[string]string)) {
	log.Debug("Route Registered: " + path)
	
	route := Route{path, method, consumes, produces, handler}
    routes = append(routes, route)
}

func init() {
	log.Debug("Initilizing router") 
    http.HandleFunc("/", handler)
}

//Handles all incomming requests and routes them to registered handlers
func handler(resp http.ResponseWriter, req *http.Request) {
	for _, route := range routes {
		reqRoute := Route{req.URL.Path, ToMethod(req.Method), nil, nil, nil}
		if route.matchRoute(reqRoute) {
			route.Handler(route, nil, nil, nil);
		}
	}
}
