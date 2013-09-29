package alert

type Check struct {
	Project string

	Alert string

	PreviousState ALERT_STATE

	CurrentState ALERT_STATE

	Changed bool

	Value float64
}
