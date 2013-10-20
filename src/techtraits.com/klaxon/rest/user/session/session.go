package session

import (
	"time"
)

const USER_KEY string = "SESSION"

type Session struct {
	UserName  string
	Projects  []string
	Expires   time.Time
	Signature string
}
