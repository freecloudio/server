package models

import "time"

type User struct {
	ID      int64     `json:"id,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`

	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`

	IsAdmin     bool  `json:"is_admin,omitempty"`
	LastSession int64 `json:"last_session,omitempty"`
}
