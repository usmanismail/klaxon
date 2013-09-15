package subscription

import (
	"appengine/datastore"
	"encoding/json"
	"net/http"
	"strings"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/rest/subscription/{project_id}", router.GET, nil, nil, getSubscriptions)
	router.Register("/rest/subscription/{project_id}/{subscription_id}", router.GET, nil, nil, getSubscription)
	router.Register("/rest/subscription/{project_id}", router.POST, []string{"application/json"}, nil, postSubscription)
}

//Get all subscriptions for a given project
func getSubscriptions(request router.Request) (int, []byte) {

	query := datastore.NewQuery(SUBSCRIPTION_KEY).Filter("Project =", request.GetPathParams()["project_id"])
	subscriptions := make([]Subscription, 0)
	_, err := query.GetAll(request.GetContext(), &subscriptions)

	if err != nil {
		log.Error("Error retriving Subscriptions: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	}

	subscriptionBytes, err := json.MarshalIndent(subscriptions, "", "	")

	if err != nil {
		log.Error("Error retriving Subscriptions: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, subscriptionBytes

}

//Create/Update an subscription for the given project
func postSubscription(request router.Request) (int, []byte) {

	var subscription Subscription
	err := json.Unmarshal(request.GetContent(), &subscription)
	if err != nil {
		log.Info("error: %v", err)
		return http.StatusBadRequest, []byte(err.Error())
	}

	subscription.Project = request.GetPathParams()["project_id"]
	_, err = datastore.Put(request.GetContext(), datastore.NewKey(request.GetContext(), SUBSCRIPTION_KEY,
		subscription.Project+"-"+subscription.Name, 0, nil), &subscription)

	if err != nil {
		log.Info("error: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, nil

}

//Get a specific subscription for a project
func getSubscription(request router.Request) (int, []byte) {
	var subscription Subscription
	err := datastore.Get(request.GetContext(), datastore.NewKey(request.GetContext(),
		SUBSCRIPTION_KEY, request.GetPathParams()["project_id"]+"-"+request.GetPathParams()["subscription_id"], 0, nil), &subscription)

	if err != nil && strings.Contains(err.Error(), "no such entity") {
		return http.StatusNotFound, []byte("Subscription not found")
	} else if err != nil {
		log.Error("Error retriving Subsciption: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		subscriptionBytes, err := json.MarshalIndent(subscription, "", "	")
		if err == nil {
			return http.StatusOK, subscriptionBytes
		} else {
			log.Info("Errror %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		}

	}
}
