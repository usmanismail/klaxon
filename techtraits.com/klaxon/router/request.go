package router

import (
	"net/http"
)

type Request struct {
	Route          Route
	PathParams     map[string]string
	HttpRequest    *http.Request
	ResponseWriter http.ResponseWriter
}
