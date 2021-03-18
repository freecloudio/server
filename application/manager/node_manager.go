package manager

import (
	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/application/storage"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/sirupsen/logrus"
)

type NodeManager interface {
	CreateUserRootFolder(authCtx *authorization.Context, userID models.UserID) *fcerror.Error
	GetNodeByPath(authCtx *authorization.Context, path string) (*models.Node, *fcerror.Error)
	GetNodeByID(authCtx *authorization.Context, nodeID models.NodeID) (*models.Node, *fcerror.Error)
	CreateNode(authCtx *authorization.Context, nodeType models.NodeType, parentNodeID models.NodeID, name string) (*models.Node, bool, *fcerror.Error)
	Close()
}

func NewNodeManager(cfg config.Config, nodePersistence persistence.NodePersistenceController, fileStorage storage.FileStorageController, managers *Managers) NodeManager {
	nodeMgr := &nodeManager{
		cfg:             cfg,
		nodePersistence: nodePersistence,
		fileStorage:     fileStorage,
		managers:        managers,
	}

	managers.Node = nodeMgr
	return nodeMgr
}

type nodeManager struct {
	cfg             config.Config
	nodePersistence persistence.NodePersistenceController
	fileStorage     storage.FileStorageController
	managers        *Managers
}

func (mgr *nodeManager) Close() {
}

func (mgr *nodeManager) CreateUserRootFolder(authCtx *authorization.Context, userID models.UserID) (fcerr *fcerror.Error) {
	fcerr = authorization.EnforceSystem(authCtx)
	if fcerr != nil {
		return
	}

	trans, fcerr := mgr.nodePersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer func() { fcerr = trans.Finish(fcerr) }()

	_, fcerr = trans.CreateUserRootFolder(userID)
	if fcerr != nil {
		logrus.WithField("userID", userID).WithError(fcerr).Error("Failed to create persistence user root folder")
		return
	}

	fcerr = mgr.fileStorage.CreateUserRootFolder(userID)
	if fcerr != nil {
		logrus.WithField("userID", userID).WithError(fcerr).Error("Failed to create file storage user root folder")
		return
	}

	return
}

func (mgr *nodeManager) CreateNode(authCtx *authorization.Context, nodeType models.NodeType, parentNodeID models.NodeID, name string) (node *models.Node, created bool, fcerr *fcerror.Error) {
	// TODO: Sanitize Name

	fcerr = authorization.EnforceUser(authCtx)
	if fcerr != nil {
		return
	}

	trans, fcerr := mgr.nodePersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer func() { fcerr = trans.Finish(fcerr) }()

	// TODO: Check if already existing

	node, created, fcerr = trans.CreateNodeByID(authCtx.User.ID, nodeType, parentNodeID, name)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create node")
		return
	}
	if !created {
		logrus.WithField("node", node).Info("File or folder already exists in persistence, don't create in storage")
		return
	}

	fcerr = mgr.fileStorage.CreateEmptyFileOrFolder(node)
	if fcerr != nil {
		logrus.WithError(fcerr).WithField("node", node).Error("Failed to create empty file or folder")
	}

	return
}

func (mgr *nodeManager) GetNodeByPath(authCtx *authorization.Context, path string) (node *models.Node, fcerr *fcerror.Error) {
	//TODO: Sanitize

	fcerr = authorization.EnforceUser(authCtx)
	if fcerr != nil {
		return
	}

	trans, fcerr := mgr.nodePersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer trans.Close()

	node, fcerr = trans.GetNodeByPath(authCtx.User.ID, path)
	if fcerr != nil && fcerr.ID != fcerror.ErrNodeNotFound {
		logrus.WithError(fcerr).WithFields(logrus.Fields{"userID": authCtx.User.ID, "path": path}).Error("Failed to get node for path")
		return
	}

	return
}

func (mgr *nodeManager) GetNodeByID(authCtx *authorization.Context, nodeID models.NodeID) (node *models.Node, fcerr *fcerror.Error) {
	fcerr = authorization.EnforceUser(authCtx)
	if fcerr != nil {
		return
	}

	trans, fcerr := mgr.nodePersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer trans.Close()

	node, fcerr = trans.GetNodeByID(authCtx.User.ID, nodeID)
	if fcerr != nil && fcerr.ID != fcerror.ErrNodeNotFound {
		logrus.WithError(fcerr).WithFields(logrus.Fields{"userID": authCtx.User.ID, "nodeID": nodeID}).Error("Failed to get node for nodeID")
		return
	}

	return
}
