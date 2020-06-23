package persistence

import (
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/domain/models"

	"github.com/sirupsen/logrus"
)

type AuthPersistenceController interface {
	StartReadTransaction() (AuthPersistenceReadTransaction, error)
	StartReadWriteTransaction() (AuthPersistenceReadWriteTransaction, error)
}

type AuthPersistenceReadTransaction interface {
	CheckToken(tokenValue models.TokenValue) (*models.Token, error)
}

type AuthPersistenceReadWriteTransaction interface {
	ReadWriteTransaction
	SaveToken(token *models.Token) error
}

var authPersistenceController map[config.PersistencePluginKey]AuthPersistenceController

func RegisterAuthPersistenceController(persistencePluginKey config.PersistencePluginKey, controller AuthPersistenceController) {
	if authPersistenceController == nil {
		authPersistenceController = make(map[config.PersistencePluginKey]AuthPersistenceController)
	}

	authPersistenceController[persistencePluginKey] = controller
	markPluginUsed(persistencePluginKey)
}

func GetAuthPersistenceController(cfg config.Config) AuthPersistenceController {
	if ctrl, ok := authPersistenceController[cfg.GetAuthPersistencePluginKey()]; ok {
		return ctrl
	}
	logrus.Fatal("No AuthPersistenceController registered")
	return nil
}
