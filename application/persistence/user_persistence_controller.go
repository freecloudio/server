package persistence

import (
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"

	"github.com/sirupsen/logrus"
)

type UserPersistenceController interface {
	StartReadTransaction() (UserPersistenceReadTransaction, *fcerror.Error)
	StartReadWriteTransaction() (UserPersistenceReadWriteTransaction, *fcerror.Error)
}

type UserPersistenceReadTransaction interface {
	GetUserByID(userID models.UserID) (*models.User, *fcerror.Error)
	GetUserByEmail(email string) (*models.User, *fcerror.Error)
}

type UserPersistenceReadWriteTransaction interface {
	ReadWriteTransaction
	UserPersistenceReadTransaction
	SaveUser(*models.User) *fcerror.Error
}

var userPersistenceController = map[config.PersistencePluginKey]UserPersistenceController{}

func RegisterUserPersistenceController(persistencePluginKey config.PersistencePluginKey, controller UserPersistenceController) {
	userPersistenceController[persistencePluginKey] = controller
	markPluginUsed(persistencePluginKey)
}

func GetUserPersistenceController(cfg config.Config) UserPersistenceController {
	if ctrl, ok := userPersistenceController[cfg.GetUserPersistencePluginKey()]; ok {
		return ctrl
	}
	logrus.Fatal("No UserPersistenceController registered")
	return nil
}
