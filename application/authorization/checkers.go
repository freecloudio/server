package authorization

import (
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
)

var errUnauthorized = fcerror.NewError(fcerror.ErrIDNotAuthorized, nil)

func EnforceSystem(ctx *Context) *fcerror.Error {
	switch ctx.Type {
	case ContextTypeSystem:
		return nil
	default:
		return errUnauthorized
	}
}

func EnforceAdmin(ctx *Context) *fcerror.Error {
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

func EnforceUser(ctx *Context) *fcerror.Error {
	switch ctx.Type {
	case ContextTypeSystem, ContextTypeUser:
		return nil
	default:
		return errUnauthorized
	}
}

func EnforceSelf(ctx *Context, targetUserID models.UserID) *fcerror.Error {
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
