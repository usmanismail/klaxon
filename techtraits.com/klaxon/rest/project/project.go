package project

import (
	"bytes"
	"encoding/json"
)

const PROJECT_KEY string = "PROJECT"

type Project interface {
	GetName() string
	GetDescription() string
	GetConfig() map[string]string
	WriteJsonToBuffer() (bytes.Buffer, error)
	GetDTO() (ProjectDTO, error)
}

func ReadProjectFromJson(projectBytes []byte) (Project, error) {
	var project ProjectStruct
	err := json.Unmarshal(projectBytes, &project)
	return project, err
}
