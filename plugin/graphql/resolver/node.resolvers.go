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
)

func (r *mutationResolver) CreateNode(ctx context.Context, input model.NodeInput) (*model.NodeCreationResult, error) {
	if input.ParentNodeIdentifier.ID == nil {
		return nil, fcerror.NewError(fcerror.ErrBadRequest, fmt.Errorf("Node creation via path not yet supported"))
	}

	authCtx := r.getAuthContext(ctx)
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

func (r *nodeResolver) Owner(ctx context.Context, obj *models.Node) (*models.User, error) {
	if r.isOnlyIDRequested(ctx) {
		return &models.User{ID: obj.OwnerID}, nil
	}
	queryResolv := &queryResolver{r.Resolver}
	return queryResolv.User(ctx, (*string)(&obj.OwnerID))
}

func (r *nodeResolver) ParentNode(ctx context.Context, obj *models.Node) (*models.Node, error) {
	if obj.ParentNodeID == nil {
		return nil, nil
	}
	if r.isOnlyIDRequested(ctx) {
		return &models.Node{ID: *obj.ParentNodeID}, nil
	}
	queryResolv := &queryResolver{r.Resolver}
	return queryResolv.Node(ctx, model.NodeIdentifierInput{ID: (*string)(obj.ParentNodeID)})
}

func (r *nodeResolver) Files(ctx context.Context, obj *models.Node) ([]*models.Node, error) {
	if obj.Type != models.NodeTypeFolder {
		return nil, nil
	}

	cacheContentID := "content" + string(obj.ID)
	contentInt := r.getObjectFromContextCache(ctx, cacheContentID)
	if contentInt != nil {
		r.logger.WithField("nodeID", obj.ID).Info("Got node content from context cache")
		return contentInt.([]*models.Node), nil
	}

	authCtx := r.getAuthContext(ctx)
	content, fcerr := r.managers.Node.ListByID(authCtx, obj.ID)
	if fcerr != nil {
		return nil, fcerr
	}

	for _, node := range content {
		r.insertObjectIntoContextCache(ctx, string(node.ID), node)
	}
	r.insertObjectIntoContextCache(ctx, cacheContentID, content)

	return content, nil
}

func (r *queryResolver) Node(ctx context.Context, input model.NodeIdentifierInput) (*models.Node, error) {
	authCtx := r.getAuthContext(ctx)

	var node *models.Node
	var fcerr *fcerror.Error
	if input.ID != nil {
		nodeInt := r.getObjectFromContextCache(ctx, *input.ID)
		if nodeInt != nil {
			r.logger.WithField("nodeID", *input.ID).Info("Got node from context cache")
			return nodeInt.(*models.Node), nil
		}

		node, fcerr = r.managers.Node.GetNodeByID(authCtx, models.NodeID(*input.ID))
	} else if input.FullPath != nil {
		node, fcerr = r.managers.Node.GetNodeByPath(authCtx, *input.FullPath)
	} else {
		return nil, fcerror.NewError(fcerror.ErrBadRequest, fmt.Errorf("Node ID or FullPath missing"))
	}

	r.insertObjectIntoContextCache(ctx, string(node.ID), node)

	if fcerr != nil {
		return nil, fcerr
	}
	return node, nil
}

// Node returns generated.NodeResolver implementation.
func (r *Resolver) Node() generated.NodeResolver { return &nodeResolver{r} }

type nodeResolver struct{ *Resolver }
