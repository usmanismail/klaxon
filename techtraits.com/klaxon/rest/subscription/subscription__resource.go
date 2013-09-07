package subscription

import (
	"techtraits.com/log"
	"techtraits.com/klaxon/router"
	"net/http"
	"net/url"
)


func init() {
	router.Register("/subscription/{project_id}", router.GET, []string{"application/json"} , nil,  getSubscriptions)
    router.Register("/subscription/{project_id}/{subscription_id}", router.GET, []string{"application/json"} , nil,  getSubscription)
    router.Register("/subscription/{project_id}", router.POST, []string{"application/json"} , nil,  postSubscription)
}

//Get all subscriptions for a given project
func getSubscriptions (route router.Route, pathParams map[string]string, queryParams url.Values,headers http.Header) {

	log.Info("Get Subscriptions");
}

//Create/Update an subscription for the given project
func postSubscription (route router.Route, pathParams map[string]string, queryParams url.Values,headers http.Header) {

	log.Info("Post Subscription");
}

//Get a specific subscription for a project
func getSubscription (route router.Route, pathParams map[string]string, queryParams url.Values,headers http.Header) {

	log.Info("Get Subscription");
}