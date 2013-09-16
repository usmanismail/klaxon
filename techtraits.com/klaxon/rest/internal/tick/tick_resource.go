package tick

import (
	"appengine/datastore"
	"appengine/taskqueue"
	"net/http"
	"techtraits.com/klaxon/rest/project"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/rest/internal/tick/", router.GET, nil, nil, getTick)
}

func getTick(request router.Request) (int, []byte) {

	query := datastore.NewQuery(project.PROJECT_KEY)
	projects := make([]project.ProjectDTO, 0)
	_, err := query.GetAll(request.GetContext(), &projects)

	if err != nil {
		log.Error("Error retriving projects: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {

		for _, project := range projects {
			task := taskqueue.NewPOSTTask("/rest/internal/check/"+project.Name, nil)
			if _, err := taskqueue.Add(request.GetContext(), task, "alertCheckQueue"); err != nil {
				log.Error("Error posting to task queue: %v", err)
			}

		}

		return http.StatusOK, nil
	}
}
