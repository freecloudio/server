package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/plugin/graphql/generated"
	"github.com/freecloudio/server/plugin/graphql/model"
	"github.com/sirupsen/logrus"
)

func (r *mutationResolver) CreateNode(ctx context.Context, input model.NodeInput) (*model.NodeCreationResult, error) {
	if input.ParentNodeIdentifier.ID == nil {
		return nil, fcerror.NewError(fcerror.ErrBadRequest, fmt.Errorf("Node creation via path not yet supported"))
	}

	authCtx := getAuthContext(ctx)
	node := &models.Node{
		ParentNodeID: (*models.NodeID)(input.ParentNodeIdentifier.ID),
		Name:         input.Name,
		Type:         input.Type,
	}

	created, fcerr := r.managers.Node.CreateNode(authCtx, node)
	if fcerr != nil {
		return nil, fcerr
	}

	return &model.NodeCreationResult{
		Created: created,
		Node:    node,
	}, nil
}

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

func (r *nodeResolver) Files(ctx context.Context, obj *models.Node) ([]*models.Node, error) {
	logrus.WithField("obj", obj).Warn("Files resolver called")
	if obj.Type != models.NodeTypeFolder {
		return nil, nil
	}
	authCtx := getAuthContext(ctx)

	content, fcerr := r.managers.Node.ListByID(authCtx, obj.ID)
	if fcerr != nil {
		return nil, fcerr
	}
	return content, nil
}

func (r *queryResolver) Node(ctx context.Context, input model.NodeIdentifierInput) (*models.Node, error) {
	authCtx := getAuthContext(ctx)

	var node *models.Node
	var fcerr *fcerror.Error
	if input.ID != nil {
		node, fcerr = r.managers.Node.GetNodeByID(authCtx, models.NodeID(*input.ID))
	} else if input.FullPath != nil {
		node, fcerr = r.managers.Node.GetNodeByPath(authCtx, *input.FullPath)
	} else {
		return nil, fcerror.NewError(fcerror.ErrBadRequest, fmt.Errorf("Node ID or FullPath missing"))
	}

	if fcerr != nil {
		return nil, fcerr
	}
	return node, nil
}

// Node returns generated.NodeResolver implementation.
func (r *Resolver) Node() generated.NodeResolver { return &nodeResolver{r} }

type nodeResolver struct{ *Resolver }
