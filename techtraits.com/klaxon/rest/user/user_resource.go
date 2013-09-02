package user

import (
	"techtraits.com/log"
	"techtraits.com/klaxon/rest/router"
	"net/http"
)


func init() {
	log.Debug("Initilizing User Resource")

    router.Register("/user/{user_id}/{transaction}", router.GET, nil,nil,  getUsers)
}

func getUsers (route router.Route, pathParams map[string]string, queryParams map[string]string,headers http.Header) {

	log.Info("Callback called");
}