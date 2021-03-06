package models

import (
	"time"
)

type UserID string

type User struct {
	ID      UserID    `json:"id" fc_neo:",unique"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" fc_neo:",unique"`
	Password  string `json:"password,omitempty"`

	IsAdmin bool `json:"is_admin"`
}

type UserUpdate struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`

	IsAdmin *bool `json:"is_admin"`
}
