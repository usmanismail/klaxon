package user

import (
	"net/http"
	"net/url"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/user/", router.GET, []string{"application/json"}, nil, getUsers)
	router.Register("/user/{user_id}/", router.GET, []string{"application/json"}, nil, getUser)
	router.Register("/user/", router.POST, []string{"application/json"}, nil, postUser)
	router.Register("/user/{user_id}/projects/{project_id}", router.GET, []string{"application/json"}, nil, getUserProjects)
}

func getUsers(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Get Users")
}

func postUser(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Post User")
}

func getUser(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Get User")
}

func getUserProjects(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Get User Projects")
}
