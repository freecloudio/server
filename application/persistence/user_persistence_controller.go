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

var userPersistenceController map[string]UserPersistenceController

func RegisterUserPersistenceController(persistencePluginKey string, controller UserPersistenceController) {
	if userPersistenceController == nil {
		userPersistenceController = make(map[string]UserPersistenceController)
	}

	userPersistenceController[persistencePluginKey] = controller
	markPluginUsed(persistencePluginKey)
}

func GetUserPersistenceController() UserPersistenceController {
	if ctrl, ok := userPersistenceController[config.GetUserPersistenceImplKey()]; ok {
		return ctrl
	}
	logrus.Fatal("No UserPersistenceController registered")
	return nil
}
