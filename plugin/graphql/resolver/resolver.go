package resolver

import (
	"context"
	"net/http"

	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/plugin/gin/keys"
	"github.com/freecloudio/server/utils"

	"github.com/99designs/gqlgen/graphql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen

const contextKeyObjectCache = "object_cache"

type contextCache map[string]interface{}

type Resolver struct {
	cfg      config.Config
	managers *manager.Managers
	logger   utils.Logger
}

func NewResolver(cfg config.Config, managers *manager.Managers) *Resolver {
	return &Resolver{cfg, managers, utils.CreateLogger(cfg.GetLoggingConfig())}
}

func (r *Resolver) getAuthContext(ctx context.Context) *authorization.Context {
	authContextInt := ctx.Value(keys.AuthContextKey)
	if authContextInt == nil {
		r.logger.Warn("AuthContext not found in gin context")
		return authorization.NewAnonymous()
	}
	authContext, ok := authContextInt.(*authorization.Context)
	if !ok {
		r.logger.Warn("AuthContext in gin context is not of correct type")
		return authorization.NewAnonymous()
	}
	return authContext
}

func ContextCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cache := contextCache{}
		r = r.WithContext(context.WithValue(r.Context(), contextKeyObjectCache, cache))
		next.ServeHTTP(w, r)
	})
}

func (r *Resolver) getObjectCache(ctx context.Context) contextCache {
	cacheInt := ctx.Value(contextKeyObjectCache)
	if cacheInt == nil {
		r.logger.Warn("Could not get object cache for context")
		return nil
	}
	cache, ok := cacheInt.(contextCache)
	if !ok {
		r.logger.WithField("cache", cacheInt).Warn("ObjectCache in context is not of correct type")
	}
	return cache
}

func (r *Resolver) getObjectFromContextCache(ctx context.Context, id string) interface{} {
	cache := r.getObjectCache(ctx)
	if cache == nil {
		return nil
	}
	return cache[id]
}

func (r *Resolver) insertObjectIntoContextCache(ctx context.Context, id string, obj interface{}) {
	cache := r.getObjectCache(ctx)
	if cache == nil {
		return
	}
	cache[id] = obj
}

func (r *Resolver) isOnlyIDRequested(ctx context.Context) bool {
	fields := graphql.CollectAllFields(ctx)
	return len(fields) == 1 && fields[0] == "id"
}
