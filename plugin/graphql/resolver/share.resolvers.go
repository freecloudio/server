package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/plugin/graphql/generated"
	"github.com/freecloudio/server/plugin/graphql/model"
)

func (r *mutationResolver) ShareNode(ctx context.Context, input model.ShareInput) (*model.NodeShareResult, error) {
	authCtx := getAuthContext(ctx)
	share := &models.Share{
		NodeID:       models.NodeID(input.NodeID),
		SharedWithID: models.UserID(input.SharedWithID),
		Mode:         input.Mode,
	}

	created, fcerr := r.managers.Share.CreateShare(authCtx, share)
	if fcerr != nil {
		return nil, fcerr
	}
	return &model.NodeShareResult{
		Created: created,
		Share:   share,
	}, nil
}

func (r *shareResolver) Node(ctx context.Context, obj *models.Share) (*models.Node, error) {
	if isOnlyIDRequested(ctx) {
		return &models.Node{ID: obj.NodeID}, nil
	}
	queryResolv := &queryResolver{r.Resolver}
	return queryResolv.Node(ctx, model.NodeIdentifierInput{ID: (*string)(&obj.NodeID)})
}

func (r *shareResolver) SharedWith(ctx context.Context, obj *models.Share) (*models.User, error) {
	if isOnlyIDRequested(ctx) {
		return &models.User{ID: obj.SharedWithID}, nil
	}
	queryResolv := &queryResolver{r.Resolver}
	return queryResolv.User(ctx, (*string)(&obj.SharedWithID))
}

// Share returns generated.ShareResolver implementation.
func (r *Resolver) Share() generated.ShareResolver { return &shareResolver{r} }

type shareResolver struct{ *Resolver }
