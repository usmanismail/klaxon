package alert

import (
	"appengine/datastore"
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

	query := datastore.NewQuery(ALERT_KEY).Filter("Project =", request.GetPathParams()["project_id"])
	alerts := make([]Alert, 0)
	_, err := query.GetAll(request.GetContext(), &alerts)

	if err != nil {
		log.Error("Error retriving alerts: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		alertBytes, err := json.MarshalIndent(alerts, "", "	")
		if err != nil {
			log.Error("Error retriving Alerts: %v", err)
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
		log.Info("error: %v", err)
		return http.StatusBadRequest, []byte(err.Error())
	}

	alert.Project = request.GetPathParams()["project_id"]
	_, err = datastore.Put(request.GetContext(), datastore.NewKey(request.GetContext(), ALERT_KEY,
		alert.Project+"-"+alert.Name, 0, nil), &alert)
	if err != nil {
		log.Info("error: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, nil

}

//Get a specific alert for a project
func getAlert(request router.Request) (int, []byte) {
	var alert Alert
	err := datastore.Get(request.GetContext(), datastore.NewKey(request.GetContext(),
		ALERT_KEY, request.GetPathParams()["project_id"]+"-"+request.GetPathParams()["alert_id"], 0, nil), &alert)

	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Error("Error retriving Alert: %v", err)
		return http.StatusNotFound, []byte("Alert not found")
	} else if err != nil {
		log.Error("Error retriving alert: %v", err)
		return http.StatusBadRequest, []byte(err.Error())
	} else {
		var alertJSON, err = json.MarshalIndent(alert, "", "	")
		if err == nil {
			return http.StatusOK, alertJSON
		} else {
			log.Info("Errror %v", err)
			return http.StatusBadRequest, []byte(err.Error())
		}
	}
}
