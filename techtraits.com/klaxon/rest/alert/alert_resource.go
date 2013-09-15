package alert

import (
	"appengine/datastore"
	"bytes"
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
func getAlerts(request router.Request) {

	query := datastore.NewQuery(ALERT_KEY).Filter("Project =", request.GetPathParams()["project_id"])
	var alerts []Alert
	_, err := query.GetAll(request.GetContext(), &alerts)

	if err != nil {
		log.Error("Error retriving alerts: %v", err)
		http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
	} else {
		var alertBytes, _ = json.Marshal(alerts)
		var respBuffer bytes.Buffer
		json.Indent(&respBuffer, alertBytes, "", "	")
		respBuffer.WriteTo(request.GetResponseWriter())
	}
}

//Create/Update an alert for the given project
func postAlert(request router.Request) {

	alert, err := ReadAlertFromJson(request.GetContent())
	if err != nil {
		log.Info("error: %v", err)
		http.Error(request.GetResponseWriter(), err.Error(), http.StatusBadRequest)
	} else {
		alert.Project = request.GetPathParams()["project_id"]
		_, err = datastore.Put(request.GetContext(), datastore.NewKey(request.GetContext(), ALERT_KEY,
			alert.Project+"-"+alert.Name, 0, nil), &alert)
		if err != nil {
			log.Info("error: %v", err)
			http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
	}
}

//Get a specific alert for a project
func getAlert(request router.Request) {
	var alert Alert
	err := datastore.Get(request.GetContext(), datastore.NewKey(request.GetContext(),
		ALERT_KEY, request.GetPathParams()["project_id"]+"-"+request.GetPathParams()["alert_id"], 0, nil), &alert)

	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Error("Error retriving Alert: %v", err)
		http.Error(request.GetResponseWriter(), "Alert not found", http.StatusNotFound)
	} else if err != nil {
		log.Error("Error retriving alert: %v", err)
		http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
	} else {
		var alertJSON, err = alert.WriteJsonToBuffer()
		if err == nil {
			alertJSON.WriteTo(request.GetResponseWriter())
		} else {
			log.Info("Errror %v", err)
			http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
	}
}
