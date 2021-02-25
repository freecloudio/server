package persistence

import (
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"

	"github.com/sirupsen/logrus"
)

type AuthPersistenceController interface {
	StartReadTransaction() (AuthPersistenceReadTransaction, *fcerror.Error)
	StartReadWriteTransaction() (AuthPersistenceReadWriteTransaction, *fcerror.Error)
}

type AuthPersistenceReadTransaction interface {
	ReadTransaction
	CheckToken(tokenValue models.TokenValue) (*models.Token, *fcerror.Error)
}

type AuthPersistenceReadWriteTransaction interface {
	ReadWriteTransaction
	AuthPersistenceReadTransaction
	SaveToken(token *models.Token) *fcerror.Error
}

var authPersistenceController = map[config.PersistencePluginKey]AuthPersistenceController{}

func RegisterAuthPersistenceController(persistencePluginKey config.PersistencePluginKey, controller AuthPersistenceController) {
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
