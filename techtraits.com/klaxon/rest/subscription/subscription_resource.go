package subscription

import (
	"appengine/datastore"
	"bytes"
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
func getSubscriptions(request router.Request) {

	query := datastore.NewQuery(SUBSCRIPTION_KEY).Filter("Project =", request.GetPathParams()["project_id"])
	var subscriptions []Subscription
	_, err := query.GetAll(request.GetContext(), &subscriptions)

	if err != nil {
		log.Error("Error retriving Subscriptions: %v", err)
		http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
	} else {
		var subscriptionBytes, _ = json.Marshal(subscriptions)
		var respBuffer bytes.Buffer
		json.Indent(&respBuffer, subscriptionBytes, "", "	")
		respBuffer.WriteTo(request.GetResponseWriter())
	}
}

//Create/Update an subscription for the given project
func postSubscription(request router.Request) {

	subscription, err := ReadSubscriptionFromJson(request.GetContent())
	if err != nil {
		log.Info("error: %v", err)
		http.Error(request.GetResponseWriter(), err.Error(), http.StatusBadRequest)
	} else {
		subscription.Project = request.GetPathParams()["project_id"]
		_, err = datastore.Put(request.GetContext(), datastore.NewKey(request.GetContext(), SUBSCRIPTION_KEY,
			subscription.Project+"-"+subscription.Name, 0, nil), &subscription)
		if err != nil {
			log.Info("error: %v", err)
			http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
	}
}

//Get a specific subscription for a project
func getSubscription(request router.Request) {
	var subscription Subscription
	err := datastore.Get(request.GetContext(), datastore.NewKey(request.GetContext(),
		SUBSCRIPTION_KEY, request.GetPathParams()["project_id"]+"-"+request.GetPathParams()["subscription_id"], 0, nil), &subscription)

	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Error("Error retriving Subsciption: %v", err)
		http.Error(request.GetResponseWriter(), "Subsciption not found", http.StatusNotFound)
	} else if err != nil {
		log.Error("Error retriving Subsciption: %v", err)
		http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
	} else {
		var subscriptionJSON, err = subscription.WriteJsonToBuffer()
		if err == nil {
			subscriptionJSON.WriteTo(request.GetResponseWriter())
		} else {
			log.Info("Errror %v", err)
			http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
		}

	}
}
