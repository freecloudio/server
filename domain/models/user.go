package models

import "time"

type UserID int64

type User struct {
	ID      UserID    `json:"id,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`

	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`

	IsAdmin bool `json:"is_admin,omitempty"`
}
