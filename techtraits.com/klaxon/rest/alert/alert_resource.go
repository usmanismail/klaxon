package alert

import (
	"net/http"
	"net/url"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/alert/{project_id}", router.GET, []string{"application/json"}, nil, getAlerts)
	router.Register("/alert/{project_id}/{alert_id}", router.GET, []string{"application/json"}, nil, getAlert)
	router.Register("/alert/{project_id}", router.POST, []string{"application/json"}, nil, postAlert)
}

//Get all alerts for a given project
func getAlerts(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Get Alerts")
}

//Create/Update an alert for the given project
func postAlert(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Post Alert")
}

//Get a specific alert for a project
func getAlert(route router.Route, pathParams map[string]string, queryParams url.Values, headers http.Header) {

	log.Info("Get Alert")
}
