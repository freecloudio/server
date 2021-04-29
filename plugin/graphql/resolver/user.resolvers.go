package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/plugin/graphql/generated"
	"github.com/freecloudio/server/plugin/graphql/model"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.UserInput) (*models.User, error) {
	newUser := &models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
	}

	_, fcerr := r.managers.User.CreateUser(newUser)
	if fcerr != nil {
		return nil, fcerr
	}

	return newUser, nil
}

func (r *queryResolver) User(ctx context.Context, userID *string) (*models.User, error) {
	authContext := getAuthContext(ctx)

	// Get own user
	if userID == nil {
		if authContext.User != nil {
			return authContext.User, nil
		}
		return nil, fcerror.NewError(fcerror.ErrUnauthorized, nil)
	}

	// Get user by ID
	user, fcerr := r.managers.User.GetUserByID(authContext, models.UserID(*userID))
	if fcerr != nil {
		return nil, fcerr
	}

	return user, nil
}

func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return string(obj.ID), nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
