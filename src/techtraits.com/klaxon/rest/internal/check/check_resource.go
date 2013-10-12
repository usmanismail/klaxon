package check

import (
	"appengine"
	"appengine/mail"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"techtraits.com/graphite"
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

func processAlerts(projectObj project.Project, alerts []alert.Alert,
	subscriptions []subscription.Subscription, context appengine.Context) (int, []byte) {

	graphiteReader, _ := graphite.MakeGraphiteReader(projectObj.GetConfig()["graphite.baseurl"],
		projectObj.GetConfig()["graphite.lookback"], context)

	alertChecks := make([]alert.Check, 0)

	for _, projectAlert := range alerts {
		//TODO Make this in to a go routine
		alertCheck := processAlert(projectAlert, context, graphiteReader, subscriptions, alertChecks)
		alertChecks = append(alertChecks, alertCheck)
	}

	alertBytes, err := json.MarshalIndent(alertChecks, "", "	")

	if err != nil {
		log.Errorf(context, "Error marshalling response %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		return http.StatusOK, alertBytes
	}

}

func processAlert(projectAlert alert.Alert, context appengine.Context,
	graphiteReader graphite.GraphiteReader, subscriptions []subscription.Subscription,
	alertChecks []alert.Check) alert.Check {

	value, err := graphiteReader.ReadValue(projectAlert.Target)
	if err != nil {
		log.Errorf(context, "Error processing Alert: %v", err)
		saveChangeIfNeeded(projectAlert, projectAlert.PreviousState != alert.UNKNOWN, alert.UNKNOWN, context)
		return alert.Check{projectAlert.Project, projectAlert.Name, projectAlert.PreviousState,
			alert.UNKNOWN, projectAlert.PreviousState != alert.UNKNOWN, 0}

	} else {
		changed, previous, current := projectAlert.CheckAlertStatusChange(value)
		saveChangeIfNeeded(projectAlert, changed, current, context)
		check := alert.Check{projectAlert.Project, projectAlert.Name, previous, current, changed, value}
		triggerSubscriptionsIfNeeded(check, subscriptions, context)
		return check
	}
}

func saveChangeIfNeeded(projectAlert alert.Alert, changed bool, current alert.ALERT_STATE, context appengine.Context) {
	if changed {
		projectAlert.PreviousState = current
		alert.SaveAlertToGAE(projectAlert, context)
	}
}

func triggerSubscriptionsIfNeeded(check alert.Check, subscriptions []subscription.Subscription,
	context appengine.Context) {
	for _, subscription := range subscriptions {
		log.Infof(context, "Firing Subscrption for check %v", subscription, check)

		//url := createConfirmationURL(r)
		msg := &mail.Message{
			Sender:  "Klaxon <0xfffffff@gmail.com>",
			To:      []string{subscription.Target},
			Subject: "Alert Triggered",
			Body:    "Alert Triggered",
		}
		if err := mail.Send(context, msg); err != nil {
			log.Errorf(context, "Couldn't send email: %v", err)
		}

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
