package project

import (
	"appengine/datastore"
	"encoding/json"
	"net/http"
	"strings"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {
	router.Register("/rest/project", router.GET, nil, nil, getProjects)
	router.Register("/rest/project/{project_id}", router.GET, nil, nil, getProject)
	router.Register("/rest/project", router.POST, []string{"application/json"}, nil, postProject)
}

//Get all projects
func getProjects(request router.Request) (int, []byte) {

	query := datastore.NewQuery(PROJECT_KEY)

	projectDTOs := make([]ProjectDTO, 0)
	_, err := query.GetAll(request.GetContext(), &projectDTOs)

	if err != nil {
		log.Error("Error retriving project: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		//Convert to Projects
		var projects []Project
		for _, projectDTO := range projectDTOs {
			var project, err = projectDTO.GetProject()
			if err == nil {
				projects = append(projects, project)
			} else {
				return http.StatusInternalServerError, []byte(err.Error())
			}
		}

		projectBytes, err := json.MarshalIndent(projects, "", "	")
		if err != nil {
			log.Error("Error retriving Projects: %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		}
		return http.StatusOK, projectBytes
	}
}

//Create/Update a project
func postProject(request router.Request) (int, []byte) {

	var project ProjectStruct
	err := json.Unmarshal(request.GetContent(), &project)
	if err != nil {
		log.Info("error: %v", err)
		return http.StatusBadRequest, []byte(err.Error())
	}

	projectDTO, err := project.GetDTO()
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())

	}

	_, err = datastore.Put(request.GetContext(), datastore.NewKey(request.GetContext(), PROJECT_KEY, project.GetName(), 0, nil), &projectDTO)
	if err != nil {
		log.Info("error: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, nil
}

//Get a specific project
func getProject(request router.Request) (int, []byte) {

	var projectDTO ProjectDTO
	err := datastore.Get(request.GetContext(), datastore.NewKey(request.GetContext(),
		PROJECT_KEY, request.GetPathParams()["project_id"], 0, nil), &projectDTO)

	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Error("Error retriving Project: %v", err)
		return http.StatusNotFound, []byte("Project Not Found")
	} else if err != nil {
		log.Error("Error retriving project: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		project, err := projectDTO.GetProject()
		if err != nil {
			log.Info("Errror %v", err)
			return http.StatusInternalServerError, []byte(err.Error())

		}

		projectJSON, err := json.MarshalIndent(project, "", "	")
		if err != nil {
			log.Info("Errror %v", err)
			return http.StatusInternalServerError, []byte(err.Error())

		}

		return http.StatusOK, projectJSON

	}

}
