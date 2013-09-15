package project

import (
	"encoding/json"
)

type ProjectStruct struct {
	//Must be unique
	Name string

	//A human readable escription of the project
	Description string

	// A key value pair of all settings related to the project
	Config map[string]string
}

func (this ProjectStruct) GetName() string {
	return this.Name
}

func (this ProjectStruct) GetDescription() string {
	return this.Description
}

func (this ProjectStruct) GetConfig() map[string]string {
	return this.Config
}

func (this ProjectStruct) GetDTO() (ProjectDTO, error) {
	var project ProjectDTO
	project.Name = this.Name
	project.Description = this.Description

	var err error
	project.Config, err = json.Marshal(this.Config)
	return project, err
}
