package alert

type Check struct {
	Project string

	Alert string

	PreviousState ALERT_STATE

	CurrentState ALERT_STATE

	Changed bool

	Value float64
}

func GetStateString(state ALERT_STATE) string {

	switch state {
	default:
		return "Unknown"
	case 1:
		return "Ok"
	case 2:
		return "Warn"
	case 3:
		return "Error"
	}
}
