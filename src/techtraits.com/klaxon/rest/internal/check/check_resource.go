package check

import (
	"appengine"
	"encoding/json"
	"net/http"
	"strings"
	"techtraits.com/graphite"
	"techtraits.com/klaxon/rest/alert"
	"techtraits.com/klaxon/rest/project"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/rest/internal/check/{project_id}", router.POST, nil, nil, checkProject)
}

func checkProject(request router.Request) (int, []byte) {

	//Check that Project exists
	projectDto, err := project.GetProjectDTOFromGAE(request.GetPathParams()["project_id"], request.GetContext())

	//TODO: Check if project is enabled
	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Errorf(request.GetContext(), "Error retriving project: %v", err)
		return http.StatusNotFound, []byte(err.Error())
	} else if err != nil {
		log.Errorf(request.GetContext(), "Error retriving project: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		projectObj, err := projectDto.GetProject()
		if err != nil {
			log.Warnf(request.GetContext(), "Error polling graphite %v ", err.Error())
			return http.StatusInternalServerError, []byte(err.Error())
		} else if projectObj.GetConfig()["graphite.baseurl"] == "" {
			log.Warnf(request.GetContext(), "Error polling graphite, property graphite.base missing for project ")
			return http.StatusInternalServerError, []byte("Property graphite.base missing for project")
		} else {
			return processAlerts(request.GetPathParams()["project_id"], projectObj.GetConfig()["graphite.baseurl"], request.GetContext())
		}
	}

}

func processAlerts(projectId string, graphiteUrl string, context appengine.Context) (int, []byte) {

	//Get Alerts for Project
	alerts, err := alert.GetAlertsFromGAE(projectId, context)
	if err != nil {
		log.Errorf(context, "Error retriving alerts: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		alertChecks := make([]alert.Check, 0)
		graphiteReader, _ := graphite.MakeGraphiteReader(graphiteUrl, context)
		for _, projectAlert := range alerts {
			alertChecks = append(alertChecks, processAlert(projectAlert, context, graphiteReader, alertChecks))
		}
		alertBytes, err := json.MarshalIndent(alertChecks, "", "	")
		if err != nil {
			log.Errorf(context, "Error marshalling response %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		}
		return http.StatusOK, alertBytes
	}
}

func processAlert(projectAlert alert.Alert, context appengine.Context, graphiteReader graphite.GraphiteReader, alertChecks []alert.Check) alert.Check {
	value, err := graphiteReader.ReadValue(projectAlert.Target)
	if err != nil {
		log.Errorf(context, "Error processing Alert: %v", err)
		saveChangeIfNeeded(projectAlert, projectAlert.PreviousState != alert.UNKNOWN, alert.UNKNOWN, context)
		return alert.Check{projectAlert.Project, projectAlert.Name, projectAlert.PreviousState, alert.UNKNOWN, projectAlert.PreviousState != alert.UNKNOWN, 0}

	} else {
		changed, previous, current := projectAlert.CheckAlertStatusChange(value)
		saveChangeIfNeeded(projectAlert, changed, current, context)
		log.Info("Send Alert %v changed", projectAlert.Name, value, previous, current, changed)
		return alert.Check{projectAlert.Project, projectAlert.Name, previous, current, changed, value}
	}
}

func saveChangeIfNeeded(projectAlert alert.Alert, changed bool, current alert.ALERT_STATE, context appengine.Context) {
	if changed {
		projectAlert.PreviousState = current
		alert.SaveAlertToGAE(projectAlert, context)
	}
}
