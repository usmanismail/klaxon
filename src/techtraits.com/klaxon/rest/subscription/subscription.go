package subscription

import (
	"appengine"
	"appengine/datastore"
	"appengine/mail"
	"fmt"
	"techtraits.com/klaxon/rest/alert"
	"techtraits.com/log"
	"time"
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

func TriggerSubscriptionsIfNeeded(check alert.Check, subscriptions []Subscription, context appengine.Context) {
	for _, subscription := range subscriptions {

		if check.Changed && check.CurrentState > check.PreviousState {
			log.Infof(context, "Firing Subscrption %v", subscription)

			subject := fmt.Sprintf("Alert %s changed state from %s to %s", check.Alert,
				alert.GetStateString(check.PreviousState), alert.GetStateString(check.CurrentState))

			message := fmt.Sprintf("Alert %s changed state from %s to %s with value %f\n Value measured at %s.\n", check.Alert,
				alert.GetStateString(check.PreviousState), alert.GetStateString(check.CurrentState),
				check.Value, time.Now().UTC().String())

			msg := &mail.Message{
				Sender:  "Klaxon <usman@techtraits.com>",
				To:      []string{subscription.Target},
				Subject: subject,
				Body:    message,
			}
			if err := mail.Send(context, msg); err != nil {
				log.Errorf(context, "Couldn't send email: %v", err)
			}
		}

	}
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
