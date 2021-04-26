package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
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
	ginCtx, fcerr := extractGinContext(ctx)
	if fcerr != nil {
		return nil, fcerr
	}
	authContext := getAuthContext(ginCtx)

	if tokenInt, ok := ginCtx.Get(AuthTokenKey); authContext.Type == authorization.ContextTypeUser && ok {
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

func (r *sessionResolver) UserID(ctx context.Context, obj *models.Session) (string, error) {
	return string(obj.UserID), nil
}

// Session returns generated.SessionResolver implementation.
func (r *Resolver) Session() generated.SessionResolver { return &sessionResolver{r} }

type sessionResolver struct{ *Resolver }
