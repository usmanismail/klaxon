package alert

import (
	"encoding/json"
	"net/http"
	"strings"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/rest/alert/{project_id}", router.GET, nil, nil, getAlerts)
	router.Register("/rest/alert/{project_id}/{alert_id}", router.GET, nil, nil, getAlert)
	router.Register("/rest/alert/{project_id}", router.POST, []string{"application/json"}, nil, postAlert)
}

//Get all alerts for a given project
func getAlerts(request router.Request) (int, []byte) {
	alerts, err := GetAlertsFromGAE(request.GetPathParams()["project_id"], request.GetContext())
	if err != nil {
		log.Errorf(request.GetContext(), "Error retriving alerts: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		alertBytes, err := json.MarshalIndent(alerts, "", "	")
		if err != nil {
			log.Errorf(request.GetContext(), "Error retriving Alerts: %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		}
		return http.StatusOK, alertBytes
	}
}

//Create/Update an alert for the given project
func postAlert(request router.Request) (int, []byte) {

	var alert Alert
	err := json.Unmarshal(request.GetContent(), &alert)
	if err != nil {
		log.Infof(request.GetContext(), "error: %v", err)
		return http.StatusBadRequest, []byte(err.Error())
	}

	//TODO Check Project Exists
	alert.Project = request.GetPathParams()["project_id"]
	err = SaveAlertToGAE(alert, request.GetContext())
	if err != nil {
		log.Infof(request.GetContext(), "error: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, nil

}

//Get a specific alert for a project
func getAlert(request router.Request) (int, []byte) {

	alert, err := GetAlertFromGAE(request.GetPathParams()["project_id"], request.GetPathParams()["alert_id"], request.GetContext())

	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Errorf(request.GetContext(), "Error retriving Alert: %v", err)
		return http.StatusNotFound, []byte("Alert not found")
	} else if err != nil {
		log.Errorf(request.GetContext(), "Error retriving alert: %v", err)
		return http.StatusBadRequest, []byte(err.Error())
	} else {
		var alertJSON, err = json.MarshalIndent(alert, "", "	")
		if err == nil {
			return http.StatusOK, alertJSON
		} else {
			log.Errorf(request.GetContext(), "Errror %v", err)
			return http.StatusBadRequest, []byte(err.Error())
		}
	}
}
