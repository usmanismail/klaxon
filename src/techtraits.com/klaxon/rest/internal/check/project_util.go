package check

import (
	"appengine"
	"errors"
	"net/http"
	"strings"
	"techtraits.com/klaxon/rest/project"
	"techtraits.com/log"
)

func getProject(projectId string, context appengine.Context) (project.Project, int, error) {
	projectDto, err := project.GetProjectDTOFromGAE(projectId, context)
	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Errorf(context, "Error retriving project: %v", err)
		return nil, http.StatusNotFound, err
	} else if err != nil {
		log.Errorf(context, "Error retriving project: %v", err)
		return nil, http.StatusInternalServerError, err
	} else {
		projectObj, err := projectDto.GetProject()
		status, err := verifyProject(projectObj, err, context)
		if err != nil {
			return nil, status, err
		} else {
			return projectObj, 0, nil
		}
	}
}

func verifyProject(projectObj project.Project, err error, context appengine.Context) (int, error) {
	if err != nil {
		log.Warnf(context, "Error polling graphite %v ", err.Error())
		return http.StatusInternalServerError, err
	} else if projectObj.GetConfig()["graphite.baseurl"] == "" {
		log.Warnf(context, "Error polling graphite, property graphite.base missing for project ")
		return http.StatusInternalServerError, errors.New("Property graphite.base missing for project")
	} else if projectObj.GetConfig()["graphite.lookback"] == "" {
		log.Warnf(context, "Error polling graphite, property graphite.lookback missing for project ")
		return http.StatusInternalServerError, errors.New("Property graphite.lookback missing for project")
	} else {
		return 0, nil
	}
}
