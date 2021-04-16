package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/plugin/graphql/generated"

	"github.com/google/uuid"
)

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	return []*models.User{{ID: models.UserID(uuid.New().String()), FirstName: "Test", LastName: "Tester"}}, nil
}

func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return string(obj.ID), nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
