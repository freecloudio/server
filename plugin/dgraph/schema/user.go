package schema

import "github.com/freecloudio/server/domain/models"

type DUser struct {
	models.User
	UID   string   `json:"uid,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

func CreateDUser(user *models.User) *DUser {
	return &DUser{
		User:  *user,
		DType: []string{"User"},
	}
}

const User = `
	first_name: string @index(hash) .
	last_name:  string @index(hash) .
	email:      string @index(hash) .
	password:   string .
	is_admin:   bool @index(bool) .

	type User {
		first_name
		last_name
		email
		password
		is_admin
		created
		updated
	}

	token:       string @index(hash) .
	valid_until: dateTime .
	for_user:    uid .

	type Token {
		token
		valid_until
		for_user
	}
`
