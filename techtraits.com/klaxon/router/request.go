package router

import (
	"net/http"
	"appengine" 
)

type RequestStruct struct {
	Route          Route
	PathParams     map[string]string
	HttpRequest    *http.Request
	ResponseWriter http.ResponseWriter
}

type Request interface {
	GetRoute() Route
	GetPathParams() map[string]string
	GetHttpRequest() *http.Request
	GetResponseWriter() http.ResponseWriter
	GetContext() appengine.Context
	GetContent() []byte
}

func (this RequestStruct) GetRoute() Route {
	return this.Route
}

func (this RequestStruct) GetPathParams() map[string]string {
	return this.PathParams
}

func (this RequestStruct) GetHttpRequest() *http.Request {
	return this.HttpRequest
}

func (this RequestStruct) GetResponseWriter() http.ResponseWriter {
	return this.ResponseWriter
}

func (this RequestStruct) GetContext() appengine.Context {
	return appengine.NewContext(this.HttpRequest)
}

func (this RequestStruct) GetContent() []byte {
	content := make([]byte, this.HttpRequest.ContentLength)
	this.HttpRequest.Body.Read(content)
	return content
}