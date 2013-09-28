package check

import (
	"net/http"
	"techtraits.com/klaxon/rest/alert"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/rest/internal/check/{project_id}", router.POST, nil, nil, getTick)
}

func getTick(request router.Request) (int, []byte) {

	//Check that Project exists

	//Get Alerts for Project
	alerts, err := alert.GetAlertsFromGAE(request.GetPathParams()["project_id"], request.GetContext())
	if err != nil {
		log.Errorf(request.GetContext(), "Error retriving alerts: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		log.Infof(request.GetContext(), "Got alerts %v", alerts)

	}

	return http.StatusOK, nil

}
