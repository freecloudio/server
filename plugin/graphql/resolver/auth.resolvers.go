package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/plugin/gin/keys"
	"github.com/freecloudio/server/plugin/graphql/generated"
	"github.com/freecloudio/server/plugin/graphql/model"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*models.Session, error) {
	token, fcerr := r.managers.Auth.Login(input.Email, input.Password)
	if fcerr != nil {
		return nil, fcerr
	}
	return token, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (*model.MutationResult, error) {
	authContext := getAuthContext(ctx)

	var fcerr *fcerror.Error
	if tokenInt := ctx.Value(keys.AuthTokenKey); authContext.Type == authorization.ContextTypeUser && tokenInt != nil {
		token := tokenInt.(models.Token)
		fcerr = r.managers.Auth.Logout(token)
	} else {
		fcerr = fcerror.NewError(fcerror.ErrUnauthorized, nil)
	}

	if fcerr != nil {
		return nil, fcerr
	}
	return &model.MutationResult{Success: true}, nil
}

func (r *sessionResolver) Token(ctx context.Context, obj *models.Session) (string, error) {
	return string(obj.Token), nil
}

func (r *sessionResolver) User(ctx context.Context, obj *models.Session) (*models.User, error) {
	if isOnlyIDRequested(ctx) {
		return &models.User{ID: obj.UserID}, nil
	}
	queryResolv := &queryResolver{r.Resolver}
	return queryResolv.User(ctx, (*string)(&obj.UserID))
}

// Session returns generated.SessionResolver implementation.
func (r *Resolver) Session() generated.SessionResolver { return &sessionResolver{r} }

type sessionResolver struct{ *Resolver }
