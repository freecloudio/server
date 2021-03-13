package authorization

import (
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
)

var errUnauthorized = fcerror.NewError(fcerror.ErrUnauthorized, nil)
var errForbidden = fcerror.NewError(fcerror.ErrUnauthorized, nil)

func EnforceSystem(ctx *Context) *fcerror.Error {
	switch ctx.Type {
	case ContextTypeSystem:
		return nil
	case ContextTypeUser:
		return errForbidden
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
		return errForbidden
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
		return errForbidden
	default:
		return errUnauthorized
	}
}
