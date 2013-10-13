package check

import (
	"appengine"
	"errors"
	"net/http"
	"strings"
	"techtraits.com/klaxon/rest/alert"
	"techtraits.com/klaxon/rest/project"
	"techtraits.com/klaxon/rest/subscription"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
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

func getProject(projectId string, context appengine.Context) (project.Project, int, error) {
	projectDto, err := project.GetProjectDTOFromGAE(projectId, context)
	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Errorf(context, "Error retriving project: %v", err)
		return nil, http.StatusNotFound, err
	} else if err != nil {
		log.Errorf(context, "Error retriving project: %v", err)
		return nil, http.StatusInternalServerError, err
	} else {
		projectObj, err := projectDto.GetProject()
		status, err := verifyProject(projectObj, err, context)
		if err != nil {
			return nil, status, err
		} else {
			return projectObj, 0, nil
		}
	}
}

func verifyProject(projectObj project.Project, err error, context appengine.Context) (int, error) {
	if err != nil {
		log.Warnf(context, "Error polling graphite %v ", err.Error())
		return http.StatusInternalServerError, err
	} else if projectObj.GetConfig()["graphite.baseurl"] == "" {
		log.Warnf(context, "Error polling graphite, property graphite.base missing for project ")
		return http.StatusInternalServerError, errors.New("Property graphite.base missing for project")
	} else if projectObj.GetConfig()["graphite.lookback"] == "" {
		log.Warnf(context, "Error polling graphite, property graphite.lookback missing for project ")
		return http.StatusInternalServerError, errors.New("Property graphite.lookback missing for project")
	} else {
		return 0, nil
	}
}

func getAlerts(projectId string, context appengine.Context) ([]alert.Alert, int, error) {
	//Get Alerts for Project
	alerts, err := alert.GetAlertsFromGAE(projectId, context)
	if err != nil {
		log.Errorf(context, "Error retriving alerts: %v", err)
		return nil, http.StatusInternalServerError, err
	} else {
		return alerts, http.StatusOK, nil
	}
}

func getSubscriptions(projectId string, context appengine.Context) ([]subscription.Subscription, int, error) {
	//Get Subscriptions for Project
	subscriptions, err := subscription.GetSubscriptionsFromGAE(projectId, context)
	if err != nil {
		log.Errorf(context, "Error retriving subscriptions: %v", err)
		return nil, http.StatusInternalServerError, err
	} else {
		return subscriptions, http.StatusOK, nil
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
