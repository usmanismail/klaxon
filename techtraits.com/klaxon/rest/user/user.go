package user

type User struct {

	// Must be unique
	UserName string

	// The set of projects to which hthis user has access
	Projects []string

	// A hash of the password for this user
	PasswordHash string

	// Is the user an admin
	isAdmin bool
}
