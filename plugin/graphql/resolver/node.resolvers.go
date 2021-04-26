package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/plugin/graphql/generated"
)

func (r *nodeResolver) ID(ctx context.Context, obj *models.Node) (string, error) {
	return string(obj.ID), nil
}

func (r *nodeResolver) MimeType(ctx context.Context, obj *models.Node) (string, error) {
	return string(obj.MimeType), nil
}

func (r *nodeResolver) OwnerID(ctx context.Context, obj *models.Node) (string, error) {
	return string(obj.OwnerID), nil
}

func (r *nodeResolver) ParentNodeID(ctx context.Context, obj *models.Node) (*string, error) {
	if obj.ParentNodeID != nil {
		str := string(*obj.ParentNodeID)
		return &str, nil
	}
	return nil, nil
}

func (r *nodeResolver) Content(ctx context.Context, obj *models.Node) ([]*models.Node, error) {
	if obj.Type != models.NodeTypeFile {
		return nil, nil
	}
	authContext := getAuthContext(ctx)

	return r.managers.Node.ListByID(authContext, obj.ID)
}

// Node returns generated.NodeResolver implementation.
func (r *Resolver) Node() generated.NodeResolver { return &nodeResolver{r} }

type nodeResolver struct{ *Resolver }
