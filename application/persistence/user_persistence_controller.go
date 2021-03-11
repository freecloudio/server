package persistence

import (
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
)

type UserPersistenceController interface {
	StartReadTransaction() (UserPersistenceReadTransaction, *fcerror.Error)
	StartReadWriteTransaction() (UserPersistenceReadWriteTransaction, *fcerror.Error)
}

type UserPersistenceReadTransaction interface {
	ReadTransaction
	CountUsers() (int64, *fcerror.Error)
	GetUserByID(userID models.UserID) (*models.User, *fcerror.Error)
	GetUserByEmail(email string) (*models.User, *fcerror.Error)
}

type UserPersistenceReadWriteTransaction interface {
	ReadWriteTransaction
	UserPersistenceReadTransaction
	SaveUser(*models.User) *fcerror.Error
	UpdateUser(*models.User) *fcerror.Error
}
