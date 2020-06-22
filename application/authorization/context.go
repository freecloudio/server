package authorization

import "github.com/freecloudio/server/domain/models"

type ContextType int

const (
	ContextTypeSystem ContextType = iota
	ContextTypeAnonymous
	ContextTypeUser
)

type Context struct {
	Type ContextType
	User *models.User
}

func NewAnonymous() *Context {
	return &Context{Type: ContextTypeAnonymous}
}

func NewSystem() *Context {
	return &Context{Type: ContextTypeSystem}
}
