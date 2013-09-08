package user

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/user/", router.GET, []string{"application/json"}, nil, getUsers)
	router.Register("/user/{user_id}/", router.GET, []string{"application/json"}, nil, getUser)
	router.Register("/user/", router.POST, []string{"application/json"}, nil, postUser)
	router.Register("/user/{user_id}/projects/{project_id}", router.GET, []string{"application/json"}, nil, getUserProjects)
}

func getUsers(request router.Request) {

}

func postUser(request router.Request) {

	context := appengine.NewContext(request.HttpRequest)
	user := User{"usman", []string{"cascade"}, "hash!!!!!!!", true}
	_, err := datastore.Put(context, datastore.NewKey(context, USER_KEY, user.UserName, 0, nil), &user)

	if err != nil {
		log.Error("Error saving user: %v", err)
		http.Error(request.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}

func getUser(request router.Request) {

	context := appengine.NewContext(request.HttpRequest)
	var user User
	err := datastore.Get(context, datastore.NewKey(context, USER_KEY, request.PathParams["user_id"], 0, nil), &user)
	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Error("Error retriving user: %v", err)
		http.Error(request.ResponseWriter, "User not found", http.StatusNotFound)
	} else if err != nil {
		log.Error("Error retriving user: %v", err)
		http.Error(request.ResponseWriter, err.Error(), http.StatusInternalServerError)
	} else {
		var userBytes, _ = json.Marshal(user)
		var respBuffer bytes.Buffer
		json.Indent(&respBuffer, userBytes, "", "	")
		respBuffer.WriteTo(request.ResponseWriter)
	}
}

func getUserProjects(request router.Request) {

}
