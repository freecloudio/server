package persistence

import (
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
)

type SharePersistenceController interface {
	StartReadTransaction() (SharePersistenceReadTransaction, *fcerror.Error)
	StartReadWriteTransaction() (SharePersistenceReadWriteTransaction, *fcerror.Error)
}

type SharePersistenceReadTransaction interface {
	ReadTransaction
	NodeContainsNestedShares(nodeID models.NodeID) (bool, *fcerror.Error)
}

type SharePersistenceReadWriteTransaction interface {
	ReadWriteTransaction
	SharePersistenceReadTransaction
	CreateShare(userID models.UserID, share *models.Share, insertName string) (bool, *fcerror.Error)
}
