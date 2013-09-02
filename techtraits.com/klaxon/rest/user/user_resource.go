package user

import (
	"techtraits.com/log"
	"techtraits.com/klaxon/rest/router"
)


func init() {
	log.Debug("Initilizing User Resource")

    router.Register("/user", router.GET, nil,nil,  getUsers)
}

func getUsers (route router.Route, pathParams map[string]string, queryParams map[string]string, headerParams map[string]string) {

	log.Info("Callback called");
}