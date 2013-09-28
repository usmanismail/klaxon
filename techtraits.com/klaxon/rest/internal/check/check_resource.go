package check

import (
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
		log.Infof(request.GetContext(), "Got alerts %v", alerts)
		graphiteReader, _ := graphite.MakeGraphiteReader("https://www.hostedgraphite.com/0c90142e/3f0f94a7-376b-4562-b174-ee3e769813b3/graphite", request.GetContext())
		for _, alert := range alerts {
			value, err := graphiteReader.ReadValue(alert.Target)
			log.Info("response %v", value, err)
		}

	}

	return http.StatusOK, nil

}
