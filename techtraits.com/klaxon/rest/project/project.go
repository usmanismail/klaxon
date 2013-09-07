package project

type Project struct {
	//Must be unique
	ProjectName string

	//A human readable escription of the project
	ProjectDescription string

	// A key value pair of all settings related to the project
	Config map[string]string
}
