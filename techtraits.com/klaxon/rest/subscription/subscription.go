package subscription

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
