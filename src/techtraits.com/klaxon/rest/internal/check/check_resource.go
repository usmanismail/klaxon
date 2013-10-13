package check

import (
	"net/http"
	"techtraits.com/klaxon/router"
)

func init() {
	router.Register("/rest/internal/check/{project_id}", router.POST, nil, nil, checkProject)
}

//TODO: Check if project is enabled
func checkProject(request router.Request) (int, []byte) {

	//Get The project
	projectId := request.GetPathParams()["project_id"]
	projectObj, getProjectStatus, getProjectError := getProject(projectId, request.GetContext())
	//Get Alerts
	alerts, getAlertStatus, getAlertsError := getAlerts(projectId, request.GetContext())
	//Get Subscriptions
	subscriptions, getSubsStatus, getSubsError := getSubscriptions(projectId, request.GetContext())

	if getProjectError == nil && getAlertsError == nil && getSubsError == nil {
		//Check Alerts
		return processAlerts(projectObj, alerts, subscriptions, request.GetContext())
	} else {
		return processError(getProjectStatus, getAlertStatus, getSubsStatus,
			getProjectError, getAlertsError, getSubsError)
	}

}

func processError(getProjectStatus int, getAlertStatus int, getSubsStatus int,
	getProjectError error, getAlertsError error, getSubsError error) (int, []byte) {

	if getProjectError != nil {
		return getProjectStatus, []byte(getProjectError.Error())
	} else if getAlertsError != nil {
		return getAlertStatus, []byte(getAlertsError.Error())
	} else if getSubsError != nil {
		return getSubsStatus, []byte(getSubsError.Error())
	} else {
		return http.StatusInternalServerError, []byte("Unknown server error")
	}
}
