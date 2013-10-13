package check

import (
	"appengine"
	"encoding/json"
	"net/http"
	"techtraits.com/graphite"
	"techtraits.com/klaxon/rest/alert"
	"techtraits.com/klaxon/rest/project"
	"techtraits.com/klaxon/rest/subscription"
	"techtraits.com/log"
)

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

	log.Infof(context, "Sending Graphite request too %s", projectAlert.Target)
	value, err := graphiteReader.ReadValue(projectAlert.Target)
	if err != nil {
		log.Errorf(context, "Error processing Alert: %v", err)
		projectAlert.SaveChangeIfNeeded(projectAlert.PreviousState != alert.UNKNOWN, alert.UNKNOWN, context)
		return alert.Check{projectAlert.Project, projectAlert.Name, projectAlert.PreviousState,
			alert.UNKNOWN, projectAlert.PreviousState != alert.UNKNOWN, 0}

	} else {
		changed, previous, current := projectAlert.CheckAlertStatusChange(value)
		projectAlert.SaveChangeIfNeeded(changed, current, context)
		check := alert.Check{projectAlert.Project, projectAlert.Name, previous, current, changed, value}
		subscription.TriggerSubscriptionsIfNeeded(check, subscriptions, context)
		return check
	}
}
