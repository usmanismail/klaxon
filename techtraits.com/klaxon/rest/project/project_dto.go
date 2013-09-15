package project

import (
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
