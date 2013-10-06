package project

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
)

type ProjectDTO struct {
	//Must be unique
	Name string

	//A human readable escription of the project
	Description string

	// A key value pair of all settings related to the project
	Config []byte
}

func (this ProjectDTO) GetProject() (Project, error) {
	var project ProjectStruct
	project.Name = this.Name
	project.Description = this.Description
	err := json.Unmarshal(this.Config, &project.Config)

	return project, err
}

func GetProjectDTOsFromGAE(context appengine.Context) ([]ProjectDTO, error) {
	projectDTOs := make([]ProjectDTO, 0)

	query := datastore.NewQuery(PROJECT_KEY)
	_, err := query.GetAll(context, &projectDTOs)
	return projectDTOs, err
}

func GetProjectDTOFromGAE(projectId string, context appengine.Context) (ProjectDTO, error) {
	var project ProjectDTO
	err := datastore.Get(context, datastore.NewKey(context,
		PROJECT_KEY, projectId, 0, nil), &project)
	return project, err
}

func SaveProjectDTOToGAE(project ProjectDTO, context appengine.Context) error {
	_, err := datastore.Put(context, datastore.NewKey(context, PROJECT_KEY,
		project.Name, 0, nil), &project)
	return err
}
