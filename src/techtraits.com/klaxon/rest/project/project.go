package project

const PROJECT_KEY string = "PROJECT"

type Project interface {
	GetName() string
	GetDescription() string
	GetConfig() map[string]string
	GetDTO() (ProjectDTO, error)
}
