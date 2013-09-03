package router

import (
   	"net/http"
   	"techtraits.com/log"
   	"net/url"
   	"strings"
   	"regexp"
)

//A set of routes configured for this application
var routes []Route

// Registers a new handler
// Will check that it does not conflict with a route that is already configured
func Register(path string, method Method, consumes []string, produces []string, 
		handler func(Route, map[string]string, url.Values, http.Header)) {
	log.Debug("Route Registered: " + path)
	
	//Trim trailing slash if it exits
	var trailingSlash , _ = regexp.MatchString(".*/$",path)
	if  trailingSlash {
		path = path[0:len(path)-1]
	}
	
	route := Route{path, method, consumes, produces, handler}
    routes = append(routes, route)
}

func init() { 
    http.HandleFunc("/", handler)
}

//Handles all incomming requests and routes them to registered handlers
func handler(resp http.ResponseWriter, req *http.Request) {
	
	reqRoute := Route{req.URL.Path, ToMethod(req.Method), strings.Split(strings.Split(req.Header.Get("Content-Type"),";")[0],","), 
			strings.Split(strings.Split(req.Header.Get("Accept"),";")[0],","), nil}
	for _, route := range routes {
		var match, pathParams = route.matchRoute(reqRoute)
		if match {			
			req.ParseForm()
			route.Handler(route, pathParams, req.Form, req.Header);
			return
		}
	}
	
	http.Error(resp, "Resource not found", http.StatusNotFound)
}
