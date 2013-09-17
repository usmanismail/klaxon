package user

import (
	"appengine/datastore"
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

func getUsers(request router.Request) (int, []byte) {

	query := datastore.NewQuery(USER_KEY)

	users := make([]User, 0)
	_, err := query.GetAll(request.GetContext(), &users)

	if err != nil {
		log.Errorf(request.GetContext(), "Error retriving user: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		//Empty out password hash before seding user
		for _, user := range users {
			go func(userObj User) {
				userObj.PasswordHash = ""
			}(user)
		}

		userBytes, err := json.MarshalIndent(users, "", "	")
		if err == nil {
			return http.StatusOK, userBytes
		} else {
			log.Errorf(request.GetContext(), "Errror %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		}
	}

}

func postUser(request router.Request) (int, []byte) {
	var user User
	err := json.Unmarshal(request.GetContent(), &user)
	if err != nil {
		log.Errorf(request.GetContext(), "error: %v", err)
		return http.StatusBadRequest, []byte(err.Error())
	}
	_, err = datastore.Put(request.GetContext(), datastore.NewKey(request.GetContext(),
		USER_KEY, user.UserName, 0, nil), &user)
	if err != nil {
		log.Errorf(request.GetContext(), "error: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, nil

}

func getUser(request router.Request) (int, []byte) {

	var user User
	err := datastore.Get(request.GetContext(), datastore.NewKey(request.GetContext(),
		USER_KEY, request.GetPathParams()["user_id"], 0, nil), &user)
	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Errorf(request.GetContext(), "Error retriving user: %v", err)
		return http.StatusInternalServerError, []byte("User not found")
	} else if err != nil {
		log.Errorf(request.GetContext(), "Error retriving user: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		//Empty out password has before seding user
		user.PasswordHash = ""
		userBytes, err := json.MarshalIndent(user, "", "	")
		if err == nil {
			return http.StatusOK, userBytes
		} else {
			log.Errorf(request.GetContext(), "Errror %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		}

	}
}
