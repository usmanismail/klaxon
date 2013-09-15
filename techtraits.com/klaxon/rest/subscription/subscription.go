package subscription

import (
	"bytes"
	"encoding/json"
)

const SUBSCRIPTION_KEY string = "SUNSCRIPTION"

type Subscription struct {

	//The unqiue name for this subscription
	Name string

	// The project to which this subscription belongs
	Project string

	// The target of the subscription, For email subscriptions
	// this will be the destination email address
	Target string
}

func ReadSubscriptionFromJson(subscriptionBytes []byte) (Subscription, error) {
	var subscription Subscription
	err := json.Unmarshal(subscriptionBytes, &subscription)
	return subscription, err
}

func (this Subscription) WriteJsonToBuffer() (bytes.Buffer, error) {
	var subscriptionBytes, err = json.Marshal(this)
	var respBuffer bytes.Buffer
	if err == nil {
		json.Indent(&respBuffer, subscriptionBytes, "", "	")
	}
	return respBuffer, err

}
