package project

import (
	"encoding/json"
	"net/http"
	"strings"
	"techtraits.com/klaxon/router"
	"techtraits.com/log"
)

func init() {

	//Get All Projects (Admin Only), All My Projects for others
	router.Register("/rest/project", router.GET, nil, nil, getAllProjects)

	//Get Project (TODO Permission if only if the user has access or is admin)
	router.Register("/rest/project/{project_id}", router.GET, nil, nil, getProject)

	//Update Project (TODO Permission if only if the user has access or is admin)
	router.Register("/rest/project", router.POST, []string{"application/json"}, nil, updateProject)
}

//Get all projects
//TODO: Check for Admin to return all
//TODO: If not admin only return current users projects
func getAllProjects(request router.Request) (int, []byte) {

	projectDTOs, err := GetProjectDTOsFromGAE(request.GetContext())

	if err != nil {
		log.Errorf(request.GetContext(), "Errorf retriving project: %v", err)
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
			log.Errorf(request.GetContext(), "Errorf retriving Projects: %v", err)
			return http.StatusInternalServerError, []byte(err.Error())
		}
		return http.StatusOK, projectBytes
	}
}

//Create/Update a project
//TODO: Check Admin/User Permissions
//TODO: Check required config present
//TODO: If create add to user project list
func updateProject(request router.Request) (int, []byte) {

	var project ProjectStruct
	err := json.Unmarshal(request.GetContent(), &project)
	if err != nil {
		log.Infof(request.GetContext(), "error: %v", err)
		return http.StatusBadRequest, []byte(err.Error())
	}

	projectDTO, err := project.GetDTO()
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())

	}
	err = SaveProjectDTOToGAE(projectDTO, request.GetContext())

	if err != nil {
		log.Infof(request.GetContext(), "error: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, nil
}

//Get a specific project
func getProject(request router.Request) (int, []byte) {

	projectDTO, err := GetProjectDTOFromGAE(request.GetPathParams()["project_id"], request.GetContext())

	if err != nil && strings.Contains(err.Error(), "no such entity") {
		log.Errorf(request.GetContext(), "Error retriving Project: %v", err)
		return http.StatusNotFound, []byte("Project Not Found")
	} else if err != nil {
		log.Errorf(request.GetContext(), "Error retriving project: %v", err)
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		project, err := projectDTO.GetProject()
		if err != nil {
			log.Infof(request.GetContext(), "Error %v", err)
			return http.StatusInternalServerError, []byte(err.Error())

		}

		projectJSON, err := json.MarshalIndent(project, "", "	")
		if err != nil {
			log.Infof(request.GetContext(), "Error %v", err)
			return http.StatusInternalServerError, []byte(err.Error())

		}

		return http.StatusOK, projectJSON

	}

}
