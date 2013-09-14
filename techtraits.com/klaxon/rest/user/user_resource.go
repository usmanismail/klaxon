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
	router.Register("/rest/user/", router.GET, nil, nil, getUsers)
	router.Register("/rest/user/{user_id}/", router.GET, nil, nil, getUser)
	router.Register("/rest/user/", router.POST, []string{"application/json"}, nil, postUser)
}

func getUsers(request router.Request) {

	context := appengine.NewContext(request.HttpRequest)
	query := datastore.NewQuery(USER_KEY)

	var users []User
	_, err := query.GetAll(context, &users)

	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Error("Error retriving user: %v", err)
		http.Error(request.ResponseWriter, "User not found", http.StatusNotFound)
	} else if err != nil {
		log.Error("Error retriving user: %v", err)
		http.Error(request.ResponseWriter, err.Error(), http.StatusInternalServerError)
	} else {
		//Empty out password has before seding user
		for _, user := range users {
			go func(userObj User) {
				userObj.PasswordHash = ""
			}(user)
		}

		var userBytes, _ = json.Marshal(users)
		var respBuffer bytes.Buffer
		json.Indent(&respBuffer, userBytes, "", "	")
		respBuffer.WriteTo(request.ResponseWriter)
	}

}

func postUser(request router.Request) {
	var user User
	context := appengine.NewContext(request.HttpRequest)
	content := make([]byte, request.HttpRequest.ContentLength)
	request.HttpRequest.Body.Read(content)
	err := json.Unmarshal(content, &user)
	if err != nil {
		log.Info("error: %v", err)
		http.Error(request.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
	_, err = datastore.Put(context, datastore.NewKey(context, USER_KEY, user.UserName, 0, nil), &user)
	if err != nil {
		log.Info("error: %v", err)
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
		//Empty out password has before seding user
		user.PasswordHash = ""
		var userBytes, _ = json.Marshal(user)
		var respBuffer bytes.Buffer
		json.Indent(&respBuffer, userBytes, "", "	")
		respBuffer.WriteTo(request.ResponseWriter)
	}
}
