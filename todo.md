* Clean up User Management
* Implement delete for all rest resources
* Implement put for all rest resources
* Validate Project exists for all resources
	* Subscription
	* Alert
	* User
* Make post for resources a smart post which does not erase old data
* Add "enabled" field to project resource 
* Add authentication
* Add permissions management
* Fire alerts to Subscriptions from the check resource
  * First implementation -- Done
  * Register proper sender
  * Format alert email with html 
  * Check for state change --Done
  * Alert if state changes to unknown
* Use go routines to make calls to graphite in parallel 
* Add usage counting for billing
* Add Acceptance Testing Project
