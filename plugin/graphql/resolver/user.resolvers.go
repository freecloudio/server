package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/plugin/graphql/generated"
)

func (r *queryResolver) User(ctx context.Context, userID *string) (user *models.User, err error) {
	ginCtx, fcerr := extractGinContext(ctx)
	if fcerr != nil {
		err = fcerr
		return
	}
	authContext := getAuthContext(ginCtx)

	// Get own user
	if userID == nil {
		if authContext.User != nil {
			return authContext.User, nil
		}
		return nil, fcerror.NewError(fcerror.ErrUnauthorized, nil)
	}

	// Get user by ID
	user, fcerr = r.managers.User.GetUserByID(authContext, models.UserID(*userID))
	if fcerr != nil {
		err = fcerr
		return
	}

	return
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
