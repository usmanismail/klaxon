package check

import (
	"appengine"
	"net/http"
	"techtraits.com/klaxon/rest/subscription"
	"techtraits.com/log"
)

func getSubscriptions(projectId string, context appengine.Context) ([]subscription.Subscription, int, error) {
	//Get Subscriptions for Project
	subscriptions, err := subscription.GetSubscriptionsFromGAE(projectId, context)
	if err != nil {
		log.Errorf(context, "Error retriving subscriptions: %v", err)
		return nil, http.StatusInternalServerError, err
	} else {
		return subscriptions, http.StatusOK, nil
	}
}
