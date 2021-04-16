package graphql

import (
	"net/http"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/plugin/graphql/generated"
	"github.com/freecloudio/server/plugin/graphql/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func GetGraphQLHandler(cfg config.Config, managers *manager.Managers) http.Handler {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))
}

func GetGraphQLPlaygroundHandler(cfg config.Config) http.Handler {
	return playground.Handler("GraphQL playground", "/graphql")
}
