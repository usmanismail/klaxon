package alert

import (
	"appengine"
	"appengine/datastore"
)

const ALERT_KEY string = "ALERT"

type ALERT_STATE int

const ( // iota is reset to 0
	UNKNOWN ALERT_STATE = iota // c0 == 0
	OK      ALERT_STATE = iota // c1 == 1
	WARN    ALERT_STATE = iota // c2 == 2
	ERROR   ALERT_STATE = iota // c3 == 3
)

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

	// The previous state of the Alert
	PreviousState ALERT_STATE `json:"-"`
}

// Checks if there is a change to the alert status. This functions returns if the status has changed as
// well as what the previous and current values of the status are.
func (this *Alert) CheckAlertStatusChange(value float64) (changed bool, previousState ALERT_STATE, currentState ALERT_STATE) {

	currentState = this.getCurrentState(value)
	previousState = this.PreviousState
	changed = this.PreviousState != currentState
	return
}

func (this *Alert) getCurrentState(value float64) ALERT_STATE {
	if this.ErrorLevel >= this.WarnLevel {
		return this.checkHighBad(value)
	} else {
		return this.checkLowBad(value)
	}
}

func (this *Alert) checkHighBad(value float64) ALERT_STATE {
	if value >= this.ErrorLevel {
		return ERROR
	} else if value >= this.WarnLevel {
		return WARN
	} else {
		return OK
	}
}

func (this *Alert) checkLowBad(value float64) ALERT_STATE {
	if value <= this.ErrorLevel {
		return ERROR
	} else if value <= this.WarnLevel {
		return WARN
	} else {
		return OK
	}
}

func (this *Alert) SaveChangeIfNeeded(changed bool, current ALERT_STATE, context appengine.Context) {
	if changed {
		this.PreviousState = current
		SaveAlertToGAE(*this, context)
	}
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

func SaveAlertToGAE(alert Alert, context appengine.Context) error {
	_, err := datastore.Put(context, datastore.NewKey(context, ALERT_KEY,
		alert.Project+"-"+alert.Name, 0, nil), &alert)
	return err
}
