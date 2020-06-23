package application

import (
	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/utils"
	"github.com/sirupsen/logrus"
)

type AuthManager interface {
	CreateNewToken(authCtx *authorization.Context, userID models.UserID) (*models.Token, error)
}

func NewAuthManager(cfg config.Config) AuthManager {
	return &authManager{
		cfg:             cfg,
		authPersistence: persistence.GetAuthPersistenceController(cfg),
	}
}

type authManager struct {
	cfg             config.Config
	authPersistence persistence.AuthPersistenceController
}

func (mgr *authManager) CreateNewToken(authCtx *authorization.Context, userID models.UserID) (token *models.Token, err error) {
	err = authorization.EnforceSelf(authCtx, userID)
	if err != nil {
		return
	}

	token = &models.Token{
		Value:      models.TokenValue(utils.GenerateRandomString(mgr.cfg.GetTokenValueLength())),
		ValidUntil: utils.GetTimeIn(mgr.cfg.GetTokenExpirationDuration()),
		UserID:     userID,
	}

	trans, err := mgr.authPersistence.StartReadWriteTransaction()
	if err != nil {
		logrus.WithError(err).Error("Failed to create transaction")
		return
	}
	err = trans.SaveToken(token)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user")
		trans.Rollback()
		return
	}
	err = trans.Commit()
	if err != nil {
		logrus.WithError(err).Error("Failed to commit transaction")
		return
	}
	return
}
