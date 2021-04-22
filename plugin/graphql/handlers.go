package graphql

import (
	"context"
	"errors"
	"net/http"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/plugin/graphql/generated"
	"github.com/freecloudio/server/plugin/graphql/resolver"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func GetGraphQLHandler(cfg config.Config, managers *manager.Managers) http.Handler {
	res := resolver.NewResolver(cfg, managers)
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: res}))
	server.SetErrorPresenter(errorPresenter)
	server.SetRecoverFunc(recoverFunc)

	return server
}

func GetGraphQLPlaygroundHandler(cfg config.Config) http.Handler {
	return playground.Handler("GraphQL playground", "/graphql")
}

func errorPresenter(ctx context.Context, err error) *gqlerror.Error {
	var fcerr *fcerror.Error
	if errors.As(err, &fcerr) {
		return &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: fcerr.Description,
			Extensions: map[string]interface{}{
				"id":       fcerr.ID,
				"file":     fcerr.File,
				"line":     fcerr.Line,
				"function": fcerr.Function,
				"cause":    fcerr.Cause,
			},
		}
	}

	return graphql.DefaultErrorPresenter(ctx, err)
}

func recoverFunc(ctx context.Context, err interface{}) error {
	return graphql.DefaultRecover(ctx, err)
}
