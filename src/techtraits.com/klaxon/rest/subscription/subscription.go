package subscription

import (
	"appengine"
	"appengine/datastore"
)

const SUBSCRIPTION_KEY string = "SUBSCRIPTION"

type Subscription struct {

	//The unqiue name for this subscription
	Name string

	// The project to which this subscription belongs
	Project string

	// The target of the subscription, For email subscriptions
	// this will be the destination email address
	Target string
}

func GetSubscriptionsFromGAE(projectId string, context appengine.Context) ([]Subscription, error) {
	query := datastore.NewQuery(SUBSCRIPTION_KEY).Filter("Project =", projectId)
	subscriptions := make([]Subscription, 0)
	_, err := query.GetAll(context, &subscriptions)
	return subscriptions, err
}

func GetSubscriptionFromGAE(projectId string, subscriptionId string, context appengine.Context) (Subscription, error) {
	var subscription Subscription
	err := datastore.Get(context, datastore.NewKey(context,
		SUBSCRIPTION_KEY, projectId+"-"+subscriptionId, 0, nil), &subscription)
	return subscription, err
}

func SaveSubscriptionToGAE(subscription Subscription, context appengine.Context) error {
	_, err := datastore.Put(context, datastore.NewKey(context, SUBSCRIPTION_KEY,
		subscription.Project+"-"+subscription.Name, 0, nil), &subscription)
	return err
}
