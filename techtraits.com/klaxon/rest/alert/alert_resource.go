package alert

import (
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/alert/{project_id}", router.GET, []string{"application/json"}, nil, getAlerts)
	router.Register("/alert/{project_id}/{alert_id}", router.GET, []string{"application/json"}, nil, getAlert)
	router.Register("/alert/{project_id}", router.POST, []string{"application/json"}, nil, postAlert)
}

//Get all alerts for a given project
func getAlerts(request router.Request) {

	log.Info("Get Alerts")
}

//Create/Update an alert for the given project
func postAlert(request router.Request) {

	log.Info("Post Alert")
}

//Get a specific alert for a project
func getAlert(request router.Request) {

	log.Info("Get Alert")
}
