package user

import (
	"appengine"
	"appengine/datastore"
	"strings"
	"techtraits.com/log"
)

const USER_KEY string = "USER"

type UserResgistration struct {

	// Must be unique
	UserName string

	// User Password
	Password string
}

type UserPublicData struct {
	// Must be unique
	UserName string

	// The set of projects to which this user has access
	Projects []string
}

type UserData struct {

	// Must be unique
	UserName string

	// The set of projects to which this user has access
	Projects []string

	//Password Salt
	PasswordSalt string

	// A hash of the password for this user
	PasswordHash string
}

//Get Users from GAE Data Store
func GetUsersFromGAE(context appengine.Context) ([]UserData, error) {
	query := datastore.NewQuery(USER_KEY)
	users := make([]UserData, 0)
	_, err := query.GetAll(context, &users)
	return users, err
}

//Get User from GAE Data Store
func GetUserFromGAE(userName string, context appengine.Context) (UserData, error) {
	var user UserData
	err := datastore.Get(context, datastore.NewKey(context,
		USER_KEY, userName, 0, nil), &user)
	return user, err
}

func IsUsernameAvailable(userName string, context appengine.Context) bool {
	var user UserData
	err := datastore.Get(context, datastore.NewKey(context,
		USER_KEY, userName, 0, nil), &user)

	if err == nil {
		return false
	} else if !strings.Contains(err.Error(), "no such entity") {
		log.Warnf(context, "Error checking for username availability %s", err.Error())
		return false
	} else {
		return true
	}
}

func SaveUserToGAE(user UserData, context appengine.Context) error {
	_, err := datastore.Put(context, datastore.NewKey(context,
		USER_KEY, user.UserName, 0, nil), &user)
	return err
}
