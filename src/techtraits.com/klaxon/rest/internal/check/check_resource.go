package check

import (
	"appengine"
	"encoding/json"
	"net/http"
	"techtraits.com/graphite"
	"techtraits.com/klaxon/rest/alert"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/rest/internal/check/{project_id}", router.POST, nil, nil, checkProject)
}

func checkProject(request router.Request) (int, []byte) {

	//Check that Project exists

	//Get Alerts for Project
	alerts, err := alert.GetAlertsFromGAE(request.GetPathParams()["project_id"], request.GetContext())
	if err != nil {
		log.Errorf(request.GetContext(), "Error retriving alerts: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		alertChecks := make([]alert.Check, 0)
		graphiteReader, _ := graphite.MakeGraphiteReader("https://www.hostedgraphite.com/0c90142e/3f0f94a7-376b-4562-b174-ee3e769813b3/graphite", request.GetContext())
		for _, projectAlert := range alerts {
			value, err := graphiteReader.ReadValue(projectAlert.Target)
			if err != nil {
				log.Errorf(request.GetContext(), "Error processing Alert: %v", err)
				saveChangeIfNeeded(projectAlert, projectAlert.PreviousState != alert.UNKNOWN, alert.UNKNOWN, request.GetContext())
				alertCheck := alert.Check{projectAlert.Project, projectAlert.Name, projectAlert.PreviousState, alert.UNKNOWN, projectAlert.PreviousState != alert.UNKNOWN, 0}
				alertChecks = append(alertChecks, alertCheck)
			} else {
				changed, previous, current := projectAlert.CheckAlertStatusChange(value)
				saveChangeIfNeeded(projectAlert, changed, current, request.GetContext())
				alertCheck := alert.Check{projectAlert.Project, projectAlert.Name, previous, current, changed, value}
				alertChecks = append(alertChecks, alertCheck)
				log.Info("Send Alert %v changed", projectAlert.Name, value, previous, current, changed)
			}
		}
		alertBytes, err := json.MarshalIndent(alertChecks, "", "	")
		if err != nil {
			log.Errorf(request.GetContext(), "Error marshalling response %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		}
		return http.StatusOK, alertBytes
	}

}

func saveChangeIfNeeded(projectAlert alert.Alert, changed bool, current alert.ALERT_STATE, context appengine.Context) {
	if changed {
		projectAlert.PreviousState = current
		alert.SaveAlertToGAE(projectAlert, context)
	}
}
