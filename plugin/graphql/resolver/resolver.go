package resolver

import (
	"context"
	"fmt"

	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	cfg      config.Config
	managers *manager.Managers
}

func NewResolver(cfg config.Config, managers *manager.Managers) *Resolver {
	return &Resolver{cfg, managers}
}

const (
	GinContextKey  = "gin_context"
	AuthContextKey = "authentication_context"
)

func extractGinContext(ctx context.Context) (ginCtx *gin.Context, fcerr *fcerror.Error) {
	ginCtxInt := ctx.Value(GinContextKey)
	if ginCtxInt == nil {
		fcerr = fcerror.NewError(fcerror.ErrInternalServerError, fmt.Errorf("gin context not found in context"))
		return
	}

	ginCtx, ok := ginCtxInt.(*gin.Context)
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrInternalServerError, fmt.Errorf("gin context has wrong type"))
		return
	}
	return
}

func getAuthContext(c *gin.Context) *authorization.Context {
	authContextInt, found := c.Get(AuthContextKey)
	if !found {
		logrus.Warn("AuthContext not found in gin context")
		return authorization.NewAnonymous()
	}
	authContext, ok := authContextInt.(*authorization.Context)
	if !ok {
		logrus.Warn("AuthContext in gin context is not of correct type")
		return authorization.NewAnonymous()
	}
	return authContext
}
