package persistence

import (
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
)

type AuthPersistenceController interface {
	StartReadTransaction() (AuthPersistenceReadTransaction, *fcerror.Error)
	StartReadWriteTransaction() (AuthPersistenceReadWriteTransaction, *fcerror.Error)
}

type AuthPersistenceReadTransaction interface {
	ReadTransaction
	GetSessionByToken(token models.Token) (*models.Session, *fcerror.Error)
}

type AuthPersistenceReadWriteTransaction interface {
	ReadWriteTransaction
	AuthPersistenceReadTransaction
	SaveSession(session *models.Session) *fcerror.Error
	DeleteSessionByToken(token models.Token) *fcerror.Error
	DeleteExpiredSessions() *fcerror.Error
}
