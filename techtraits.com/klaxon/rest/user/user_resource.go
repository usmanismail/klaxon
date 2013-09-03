package user

import (
	"techtraits.com/log"
	"techtraits.com/klaxon/rest/router"
	"net/http"
	"net/url"
)


func init() {
	router.Register("/user/", router.GET, []string{"application/json"} , nil,  getUsers)
    router.Register("/user/{user_id}/", router.GET, []string{"application/json"} , nil,  getUser)
    router.Register("/user/", router.POST, []string{"application/json"} , nil,  postUser)
}

func getUsers (route router.Route, pathParams map[string]string, queryParams url.Values,headers http.Header) {

	log.Info("Get Users");
}

func postUser (route router.Route, pathParams map[string]string, queryParams url.Values,headers http.Header) {

	log.Info("Post User");
}

func getUser (route router.Route, pathParams map[string]string, queryParams url.Values,headers http.Header) {

	log.Info("Get User");
}