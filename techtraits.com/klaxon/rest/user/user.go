package user

const USER_KEY string = "USER"

type User struct {

	// Must be unique
	UserName string

	// The set of projects to which hthis user has access
	Projects []string

	// A hash of the password for this user
	// Tagged to make it not serialize in responses
	PasswordHash string `json:"-"`

	// Is the user an admin
	IsAdmin bool
}
