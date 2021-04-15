package manager

import (
	"errors"
	"io"

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
	ListByID(authCtx *authorization.Context, nodeID models.NodeID) ([]*models.Node, *fcerror.Error)
	CreateNode(authCtx *authorization.Context, node *models.Node) (bool, *models.Node, *fcerror.Error)
	UploadFile(authCtx *authorization.Context, node *models.Node, uploadFilePath string) (bool, *models.Node, *fcerror.Error)
	UploadFileByID(authCtx *authorization.Context, nodeID models.NodeID, uploadFilePath string) *fcerror.Error
	DownloadFile(authCtx *authorization.Context, nodeID models.NodeID) (*models.Node, io.ReadCloser, int64, *fcerror.Error)
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

func (mgr *nodeManager) CreateNode(authCtx *authorization.Context, node *models.Node) (created bool, returnNode *models.Node, fcerr *fcerror.Error) {
	// TODO: Sanitize Name

	fcerr = authorization.EnforceUser(authCtx)
	if fcerr != nil {
		return
	}

	if node.ParentNodeID == nil {
		fcerr = fcerror.NewError(fcerror.ErrBadRequest, errors.New("ParentNodeID is missing for node creation"))
		return
	}

	trans, fcerr := mgr.nodePersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	defer func() { fcerr = trans.Finish(fcerr) }()

	returnNode, created, fcerr = trans.CreateNodeByID(authCtx.User.ID, node.Type, *node.ParentNodeID, node.Name)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create node")
		return
	}
	if !created {
		logrus.WithField("node", node).Info("File or folder already exists in persistence, don't create in storage")
		return
	}

	fcerr = mgr.fileStorage.CreateEmptyFileOrFolder(returnNode)
	if fcerr != nil {
		logrus.WithError(fcerr).WithField("node", returnNode).Error("Failed to create empty file or folder")
	}

	return
}

func (mgr *nodeManager) UploadFile(authCtx *authorization.Context, node *models.Node, uploadFilePath string) (created bool, returnNode *models.Node, fcerr *fcerror.Error) {
	node.Type = models.NodeTypeFile

	created, returnNode, fcerr = mgr.CreateNode(authCtx, node)
	if fcerr != nil {
		return
	}

	fcerr = mgr.UploadFileByID(authCtx, node.ID, uploadFilePath)
	if fcerr != nil {
		return
	}
	return
}

func (mgr *nodeManager) UploadFileByID(authCtx *authorization.Context, nodeID models.NodeID, uploadFilePath string) (fcerr *fcerror.Error) {
	node, fcerr := mgr.GetNodeByID(authCtx, nodeID)
	if fcerr != nil {
		return
	}

	fcerr = mgr.fileStorage.CopyFileFromUpload(node, uploadFilePath)
	if fcerr != nil {
		logrus.WithError(fcerr).WithField("node", node).Error("Failed to copy file from upload")
		return
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

	node, fcerr = trans.GetNodeByPath(authCtx.User.ID, path, models.ShareModeRead)
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

	node, fcerr = trans.GetNodeByID(authCtx.User.ID, nodeID, models.ShareModeRead)
	if fcerr != nil && fcerr.ID != fcerror.ErrNodeNotFound {
		logrus.WithError(fcerr).WithFields(logrus.Fields{"userID": authCtx.User.ID, "nodeID": nodeID}).Error("Failed to get node for nodeID")
		return
	}

	return
}

func (mgr *nodeManager) ListByID(authCtx *authorization.Context, nodeID models.NodeID) (node []*models.Node, fcerr *fcerror.Error) {
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

	node, fcerr = trans.ListByID(authCtx.User.ID, nodeID, models.ShareModeRead)
	if fcerr != nil && fcerr.ID != fcerror.ErrNodeNotFound {
		logrus.WithError(fcerr).WithFields(logrus.Fields{"userID": authCtx.User.ID, "nodeID": nodeID}).Error("Failed to get content for nodeID")
		return
	}

	return
}

func (mgr *nodeManager) DownloadFile(authCtx *authorization.Context, nodeID models.NodeID) (node *models.Node, reader io.ReadCloser, size int64, fcerr *fcerror.Error) {
	fcerr = authorization.EnforceUser(authCtx)
	if fcerr != nil {
		return
	}

	node, fcerr = mgr.GetNodeByID(authCtx, nodeID)
	if fcerr != nil {
		return
	}

	// TODO: Support folder download
	if node.Type == models.NodeTypeFolder {
		fcerr = fcerror.NewError(fcerror.ErrUnknown, errors.New("Downloading folder not yet supported"))
		return
	}

	// TODO: Use node size here after it is filled? => Direct FS changes will break download
	reader, size, fcerr = mgr.fileStorage.DownloadFile(node)
	if fcerr != nil {
		logrus.WithError(fcerr).WithField("node", node).Error("Failed to download file")
		return
	}
	return
}
