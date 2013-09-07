package subscription

type Subscription struct {

	// Must be unique
	SubscriptionDescription string

	// The project to which this subscription belongs
	Project string

	// The target of the subscription, For email subscriptions this will be the destination email address
	Target string
}
