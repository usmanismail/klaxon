package project

import (
	"net/http"
	"net/url"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/project", router.GET, []string{"application/json"}, nil, getProjects)
	router.Register("/project/{project_id}", router.GET, []string{"application/json"}, nil, getProjects)
	router.Register("/project", router.POST, []string{"application/json"}, nil, postProject)
}

//Get all projects
func getProjects(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Get Project")
}

//Create/Update a project
func postProject(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Post Project")
}

//Get a specific project
func getProject(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Get Project")
}
