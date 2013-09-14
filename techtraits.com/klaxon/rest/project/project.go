package project

import (
	"bytes"
	"encoding/json"
)

const PROJECT_KEY string = "PROJECT"

type ProjectDTO struct {
	//Must be unique
	Name string

	//A human readable escription of the project
	Description string

	// A key value pair of all settings related to the project
	Config []byte
}

type ProjectStruct struct {
	//Must be unique
	Name string

	//A human readable escription of the project
	Description string

	// A key value pair of all settings related to the project
	Config map[string]string
}

type Project interface {
	GetName() string
	GetDescription() string
	GetConfig() map[string]string
	WriteJsonToBuffer() (bytes.Buffer, error) 
	GetDTO() (ProjectDTO, error)
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

func (this ProjectDTO) GetProject() (Project, error) {
	var project ProjectStruct
	project.Name = this.Name
	project.Description = this.Description
	err := json.Unmarshal(this.Config, &project.Config)

	return project, err
}

func (this ProjectStruct) WriteJsonToBuffer() (bytes.Buffer, error) {
	var projectBytes, err = json.Marshal(this)
	var respBuffer bytes.Buffer
	if err == nil {
		json.Indent(&respBuffer, projectBytes, "", "	")
	}
	return respBuffer, err

}

func ReadProjectFromJson(projectBytes []byte) (Project, error) {
	var project ProjectStruct
	err := json.Unmarshal(projectBytes, &project)
	return project, err
}
