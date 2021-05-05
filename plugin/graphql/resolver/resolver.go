package resolver

import (
	"context"

	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/plugin/gin/keys"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sirupsen/logrus"
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

func getAuthContext(ctx context.Context) *authorization.Context {
	authContextInt := ctx.Value(keys.AuthContextKey)
	if authContextInt == nil {
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

func isOnlyIDRequested(ctx context.Context) bool {
	fields := graphql.CollectAllFields(ctx)
	return len(fields) == 1 && fields[0] == "id"
}
