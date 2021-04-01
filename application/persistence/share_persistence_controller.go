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
}

type SharePersistenceReadWriteTransaction interface {
	ReadWriteTransaction
	SharePersistenceReadTransaction
	CreateShare(userID models.UserID, share *models.Share) (bool, *fcerror.Error)
}
