package subscription

import (
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

	log.Info("Get Subscriptions")
}

//Create/Update an subscription for the given project
func postSubscription(request router.Request) {

	log.Info("Post Subscription")
}

//Get a specific subscription for a project
func getSubscription(request router.Request) {

	log.Info("Get Subscription")
}
