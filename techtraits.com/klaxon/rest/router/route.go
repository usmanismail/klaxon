package router

import (
	"techtraits.com/log"
	"strings"
	"regexp"
	"net/http"
	"net/url"
)

type Method int;
const (  // iota is reset to 0
        GET Method = iota  // c0 == 0
        POST Method = iota  // c1 == 1
        PUT Method = iota  // c2 == 2
		HEAD Method = iota  // c2 == 3
		DELETE Method = iota  // c2 == 4        
)

//The definition of a route
type Route struct {
	// The path which will trigger this route 
	// /hello/world maps to exactly /hello/world
	// /hello/{name} maps to /hello/ANY_STRING and returns the value in {name} as a path param 
    Path string
    
    // One of GET, PUT, POST, HEAD, DELETE
    // The request only triggers a route if it has the corresponding HTTP method type
    Method Method
    
    // The Http Content-Type header of incoming request must be set to this to match. 
    // Default value is nil which matches any request
    Consumes []string
    
    // The Http Accept header of incoming request must be set to this to match. 
    // Default value is nil which matches any request
    Produces []string
    
    //The Handler function to be called when a route is matched.
    // It returns:
    //	The Matched route
    //	A map of path parameters,
    //	A map of query string paramters
    // 	A map of headers  
    Handler func(Route, map[string]string, url.Values, http.Header)
    
    
}
//Check if the given route and the current route match. 
func (this *Route) matchConsumes(route Route)bool {
	if this.Consumes == nil {
		return true;
	} else {
		for _, thisConsumes := range this.Consumes { 
			for _, routeConsumes := range route.Consumes {
				if thisConsumes == routeConsumes {
					return true;
				}
			}
		}
		
		return false;
	}
}

//Check if the given route and the current route match. 
func (this *Route) matchProduces(route Route)bool {
	if this.Produces == nil {
		return true;
	} else {
		for _, thisProduces := range this.Produces {
			for _, routeProduces := range route.Produces {
				if thisProduces == routeProduces {
					return true;
				}
			}
		}
		return false;
	}
}

//Check if the given route and the current route match. 
func (this *Route) parseUri(route Route) (match bool, pathParams map[string]string) {
	//Trim trailing slash if it exits
	var trailingSlash , _ = regexp.MatchString(".*/$",route.Path)
	if  trailingSlash {
		route.Path = route.Path[0:len(route.Path)-1]
	}
	
	pathParams = make(map[string]string)
	var thisPathTokens = strings.Split(this.Path[1:], "/")
	var routePathTokens = strings.Split(route.Path[1:], "/")
	
	match = true	
	
	if len(thisPathTokens) != len(routePathTokens) {
		match = false
	} else {
		for i :=0; i < len(thisPathTokens) ; i++  {
			var regexMatch, _ = regexp.MatchString("^{.*}$",thisPathTokens[i])
			if regexMatch {
				pathParams[thisPathTokens[i][1:len(thisPathTokens[i])-1]]=routePathTokens[i]	
			} else if thisPathTokens[i] == routePathTokens[i] {
			} else {
				match = false
				return
			}	
		}
	}
	return
}



//Check if the given route and the current route match. 
func (this *Route) matchRoute(route Route) (bool, map[string]string) {
	var matchUri, pathParams = this.parseUri(route)
	return ((this.Method == route.Method) && this.matchConsumes(route) && this.matchProduces(route) && matchUri), pathParams
	
}

//Converts a string representation of the method type into the internal constant
func ToMethod(method string) Method {
	switch method {
    case "GET":
        return GET
    case "POST":
        return POST
	case "PUT":
        return PUT
	case "HEAD":
        return HEAD
	case "DELETE":
        return DELETE                        
    default:
        log.Error("Unable to decode Http method " + method)
        return GET
    } 
}