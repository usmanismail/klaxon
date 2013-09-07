package alert

type Alert struct {
	//Must be unique within project
	AlertName string

	//The unqiue Project Name that this alert is for
	Project string

	//A human readable escription of the alert
	AlertDescription string

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
