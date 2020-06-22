package authorization

import (
	"errors"

	"github.com/freecloudio/server/domain/models"
)

var errUnauthorized = errors.New("Unauthorized")

func EnforceSystem(ctx *Context) error {
	switch ctx.Type {
	case ContextTypeSystem:
		return nil
	default:
		return errUnauthorized
	}
}

func EnforceAdmin(ctx *Context) error {
	switch ctx.Type {
	case ContextTypeSystem:
		return nil
	case ContextTypeUser:
		if ctx.User.IsAdmin {
			return nil
		}
		fallthrough
	default:
		return errUnauthorized
	}
}

func EnforceUser(ctx *Context) error {
	switch ctx.Type {
	case ContextTypeSystem, ContextTypeUser:
		return nil
	default:
		return errUnauthorized
	}
}

func EnforceSelf(ctx *Context, targetUserID models.UserID) error {
	switch ctx.Type {
	case ContextTypeSystem:
		return nil
	case ContextTypeUser:
		if ctx.User.IsAdmin || ctx.User.ID == targetUserID {
			return nil
		}
		fallthrough
	default:
		return errUnauthorized
	}
}

// EnforceOwner -> ctx with file, EnforceHasAccess -> ctx with file+bool flag
