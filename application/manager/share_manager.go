package manager

import (
	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/sirupsen/logrus"
)

type ShareManager interface {
	ShareNodeByID(authCtx *authorization.Context, share *models.Share) (bool, *fcerror.Error)
	Close()
}

func NewShareManager(
	cfg config.Config,
	sharePersistence persistence.SharePersistenceController,
	nodePersistence persistence.NodePersistenceController,
	managers *Managers,
) ShareManager {
	shareMgr := &shareManager{
		cfg:              cfg,
		sharePersistence: sharePersistence,
		nodePersistence:  nodePersistence,
		managers:         managers,
	}

	managers.Share = shareMgr
	return shareMgr
}

type shareManager struct {
	cfg              config.Config
	sharePersistence persistence.SharePersistenceController
	nodePersistence  persistence.NodePersistenceController
	managers         *Managers
}

func (mgr *shareManager) Close() {
}

func (mgr *shareManager) ShareNodeByID(authCtx *authorization.Context, share *models.Share) (created bool, fcerr *fcerror.Error) {
	fcerr = authorization.EnforceUser(authCtx)
	if fcerr != nil {
		return
	}

	nodeTrans, fcerr := mgr.nodePersistence.StartReadTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}

	shareNode, fcerr := nodeTrans.GetNodeByID(authCtx.User.ID, share.NodeID, false)
	if fcerr != nil || fcerr.ID != fcerror.ErrNodeNotFound {
		logrus.WithError(fcerr).WithField("share", share).Error("Failed to get shareNode")
		return
	} else if fcerr != nil {
		fcerr = fcerror.NewError(fcerror.ErrShareIsInsideShare, fcerr)
		return
	}

	fcerr = nodeTrans.Close()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to close nodeTrans for ShareNodeByID - ignore for now")
		fcerr = nil
	}

	shareTrans, fcerr := mgr.sharePersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}

	containsShares, fcerr := shareTrans.NodeContainsNestedShares(share.NodeID)
	if fcerr != nil {
		logrus.WithError(fcerr).WithField("share", share).Error("Failed to get whether node contains nested shares")
		return
	} else if containsShares {
		fcerr = fcerror.NewError(fcerror.ErrShareContainsOtherShares, nil)
		return
	}

	created, fcerr = shareTrans.CreateShare(authCtx.User.ID, share, shareNode.Name)
	if fcerr != nil {
		logrus.WithError(fcerr).WithField("share", share).Error("Failed to get whether node contains nested shares")
		return
	}
	return
}
