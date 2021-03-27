package authorization

import (
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
)

func EnforceSystem(ctx *Context) *fcerror.Error {
	switch ctx.Type {
	case ContextTypeSystem:
		return nil
	case ContextTypeUser:
		return fcerror.NewErrorSkipFunc(fcerror.ErrForbidden, nil)
	default:
		return fcerror.NewErrorSkipFunc(fcerror.ErrUnauthorized, nil)
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
		return fcerror.NewErrorSkipFunc(fcerror.ErrForbidden, nil)
	default:
		return fcerror.NewErrorSkipFunc(fcerror.ErrUnauthorized, nil)
	}
}

func EnforceUser(ctx *Context) *fcerror.Error {
	switch ctx.Type {
	case ContextTypeSystem, ContextTypeUser:
		return nil
	default:
		return fcerror.NewErrorSkipFunc(fcerror.ErrUnauthorized, nil)
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
		return fcerror.NewErrorSkipFunc(fcerror.ErrForbidden, nil)
	default:
		return fcerror.NewErrorSkipFunc(fcerror.ErrUnauthorized, nil)
	}
}
