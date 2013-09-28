package alert

import (
	"appengine"
	"appengine/datastore"
)

const ALERT_KEY string = "ALERT"

type Alert struct {
	//Must be unique within project
	Name string

	//The unqiue Project Name that this alert is for
	Project string `json:"-"`

	//A human readable escription of the alert
	Description string

	//The target string to fetch the data
	Target string

	// The level at which to change alert status to a warning
	// Note:
	// if ErrorLevel is higher than WarnLevel than alert is fired when value exceedes ErrorLevel
	// If ErrorLevel is lower than warnLevel than alert is fired when value goes below ErrorLevel
	WarnLevel float64

	// The level at which to change alert status to an error and send alert to subscriptions
	// Note:
	// if ErrorLevel is higher than WarnLevel than alert is fired when value exceedes ErrorLevel
	// If ErrorLevel is lower than warnLevel than alert is fired when value goes below ErrorLevel
	ErrorLevel float64
}

func GetAlertsFromGAE(projectId string, context appengine.Context) ([]Alert, error) {
	query := datastore.NewQuery(ALERT_KEY).Filter("Project =", projectId)
	alerts := make([]Alert, 0)
	_, err := query.GetAll(context, &alerts)
	return alerts, err
}

func GetAlertFromGAE(projectId string, alertId string, context appengine.Context) (Alert, error) {
	var alert Alert
	err := datastore.Get(context, datastore.NewKey(context,
		ALERT_KEY, projectId+"-"+alertId, 0, nil), &alert)
	return alert, err
}

func PostAlertToGAE(alert Alert, context appengine.Context) error {
	_, err := datastore.Put(context, datastore.NewKey(context, ALERT_KEY,
		alert.Project+"-"+alert.Name, 0, nil), &alert)
	return err
}
