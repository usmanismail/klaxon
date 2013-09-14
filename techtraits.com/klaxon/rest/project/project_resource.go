package project

import (
	"appengine/datastore"
	"bytes"
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
func getProjects(request router.Request) {

	query := datastore.NewQuery(PROJECT_KEY)

	var projectDTOs []ProjectDTO
	_, err := query.GetAll(request.GetContext(), &projectDTOs)

	if err != nil {
		log.Error("Error retriving user: %v", err)
		http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
	} else {
		//Convert to Projects
		var projects []Project
		for _, projectDTO := range projectDTOs {
			var project, err = projectDTO.GetProject()
			if err == nil {
				projects = append(projects, project)
			} else {
				http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
				return
			}
		}

		var projectBytes, _ = json.Marshal(projects)
		var respBuffer bytes.Buffer
		json.Indent(&respBuffer, projectBytes, "", "	")
		respBuffer.WriteTo(request.GetResponseWriter())
	}
}

//Create/Update a project
func postProject(request router.Request) {

	project, err := ReadProjectFromJson(request.GetContent())
	if err != nil {
		log.Info("error: %v", err)
		http.Error(request.GetResponseWriter(), err.Error(), http.StatusBadRequest)
	} else {
		var projectDTO, err = project.GetDTO()
		if err == nil {
			_, err = datastore.Put(request.GetContext(), datastore.NewKey(request.GetContext(), PROJECT_KEY, project.GetName(), 0, nil), &projectDTO)
			if err != nil {
				log.Info("error: %v", err)
				http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
	}
}

//Get a specific project
func getProject(request router.Request) {

	var projectDTO ProjectDTO
	err := datastore.Get(request.GetContext(), datastore.NewKey(request.GetContext(),
		PROJECT_KEY, request.GetPathParams()["project_id"], 0, nil), &projectDTO)

	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Error("Error retriving Project: %v", err)
		http.Error(request.GetResponseWriter(), "Project not found", http.StatusNotFound)
	} else if err != nil {
		log.Error("Error retriving project: %v", err)
		http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
	} else {
		var project, err = projectDTO.GetProject()
		if err == nil {
			var projectJSON, err = project.WriteJsonToBuffer()
			if err == nil {
				projectJSON.WriteTo(request.GetResponseWriter())
			} else {
				log.Info("Errror %v", err)
				http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
			}
		} else {
			log.Info("Errror %v", err)
			http.Error(request.GetResponseWriter(), err.Error(), http.StatusInternalServerError)
		}
	}

}
