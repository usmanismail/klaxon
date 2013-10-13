klaxon
======
klaxon is an alerting dashboard for [Graphite](http://graphite.wikidot.com/) written in [Go](|http://golang.org) and designed to run on [Google App Engine](https://developers.google.com/appengine/). Klaxon is modeled after the [Seyren](https://github.com/scobal/seyren) framework but provides multitenancy as well as authentication.

## Status
We are currently in development, we have skeleton servlet handlers setup as well as our data objects and we hope to have the rest only implementation done fairly soon. Following which we will concentrate on the user interface. 

The tip of this code base is deployed to https://go-klaxon.appspot.com/ 

## Setup Instructions
#### Setup GAE go SDK
* [Download Google App Engine Go SDK](https://developers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go)
* Unzip the downloaded file to somewhere like /usr/local/
* Add the unzipped folder to your path

#### Build Klaxon
	git clone git@github.com:usmanismail/klaxon.git
	dev_appserver.py ./klaxon/src

#### Deploy Klaxon
	#Please replace the application tag in app.yaml with your own Application ID
	cd klaxon
	scripts/fixFormat.sh
	appcfg.py update ./src


## Contributors 
[Usman Ismail](http://techtraits.com/usman.html)