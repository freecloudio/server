package persistence

import (
	"github.com/freecloudio/server/config"
	"github.com/freecloudio/server/domain/models"
	"github.com/sirupsen/logrus"
)

type UserPersistenceController interface {
	StartTransaction() (UserPersistenceTransaction, error)
}

type UserPersistenceTransaction interface {
	Transaction
	SaveUser(*models.User) error
}

var userPersistenceController map[config.PersistencePluginKey]UserPersistenceController

func RegisterUserPersistenceController(persistencePluginKey config.PersistencePluginKey, controller UserPersistenceController) {
	if userPersistenceController == nil {
		userPersistenceController = make(map[config.PersistencePluginKey]UserPersistenceController)
	}

	userPersistenceController[persistencePluginKey] = controller
	markPluginUsed(persistencePluginKey)
}

func GetUserPersistenceController() UserPersistenceController {
	if ctrl, ok := userPersistenceController[config.GetUserPersistencePluginKey()]; ok {
		return ctrl
	}
	logrus.Fatal("No UserPersistenceController registered")
	return nil
}
