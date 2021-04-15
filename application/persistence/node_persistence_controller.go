package persistence

import (
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
)

type NodePersistenceController interface {
	StartReadTransaction() (NodePersistenceReadTransaction, *fcerror.Error)
	StartReadWriteTransaction() (NodePersistenceReadWriteTransaction, *fcerror.Error)
}

type NodePersistenceReadTransaction interface {
	ReadTransaction
	GetNodeByPath(userID models.UserID, path string, includedShareMode models.ShareMode) (*models.Node, *fcerror.Error)
	GetNodeByID(userID models.UserID, nodeID models.NodeID, includedShareMode models.ShareMode) (*models.Node, *fcerror.Error)
	ListByID(userID models.UserID, nodeID models.NodeID, includedShareMode models.ShareMode) ([]*models.Node, *fcerror.Error)
}

type NodePersistenceReadWriteTransaction interface {
	ReadWriteTransaction
	NodePersistenceReadTransaction
	CreateUserRootFolder(userID models.UserID) (bool, *fcerror.Error)
	CreateNodeByID(userID models.UserID, nodeType models.NodeType, parentNodeID models.NodeID, name string) (*models.Node, bool, *fcerror.Error)
}
