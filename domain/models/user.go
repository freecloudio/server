package models

import (
	"time"
)

type UserID int64

type User struct {
	ID      UserID    `json:"id" fc_neo:"-"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`

	IsAdmin bool `json:"is_admin"`
}
