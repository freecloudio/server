package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/freecloudio/server/plugin/graphql/generated"
	"github.com/freecloudio/server/plugin/graphql/model"
)

func (r *queryResolver) Health(ctx context.Context) (*model.MutationResult, error) {
	return &model.MutationResult{Success: true}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
