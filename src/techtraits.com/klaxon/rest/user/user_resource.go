package user

import (
	"encoding/json"
	"net/http"
	"techtraits.com/bcrypt"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	//Get All Users (Admin Only)
	router.Register("/rest/user/", router.GET, nil, nil, getUsers)

	//Register User
	router.Register("/rest/user/", router.PUT, []string{"application/json"}, nil, registerUser)

	//Get One User (Admin Only)
	//router.Register("/rest/user/{user_id}/", router.GET, nil, nil, getUser)

	//Update User (Admin Only)
	//router.Register("/rest/user/", router.PUT, []string{"application/json"}, nil, updateUser)

	// Get Me  --TODO
	// Change Password --TODO
}

func getUsers(request router.Request) (int, []byte) {
	users, err := GetUsersFromGAE(request.GetContext())

	if err != nil {
		log.Errorf(request.GetContext(), "Error retriving user: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		userBytes, err := json.MarshalIndent(users, "", "	")
		if err == nil {
			return http.StatusOK, userBytes
		} else {
			log.Errorf(request.GetContext(), "Errror %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		}
	}
}

func registerUser(request router.Request) (int, []byte) {
	var registration UserResgistration
	err := json.Unmarshal(request.GetContent(), &registration)

	//Check for correct desrialization
	if err != nil {
		log.Errorf(request.GetContext(), "error: %v", err)
		return http.StatusBadRequest, []byte(err.Error())
	}

	//check if username already exists
	if IsUsernameAvailable(registration.UserName, request.GetContext()) {

		// generate a random salt with default rounds of complexity
		salt, _ := bcrypt.Salt()

		// hash and verify a password with random salt
		hash, _ := bcrypt.Hash(registration.Password)

		var user = UserData{registration.UserName, nil, salt, hash}
		err = SaveUserToGAE(user, request.GetContext())
		if err != nil {
			log.Errorf(request.GetContext(), "error: %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		} else {
			return http.StatusOK, nil
		}
	} else {
		return http.StatusConflict, []byte("Username not avaiable")
	}
}

/*
func updateUser(request router.Request) (int, []byte) {
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
*/
