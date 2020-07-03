package application

import (
	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"
	"github.com/sirupsen/logrus"
)

type AuthManager interface {
	CreateNewToken(authCtx *authorization.Context, userID models.UserID) (*models.Token, *fcerror.Error)
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

func (mgr *authManager) CreateNewToken(authCtx *authorization.Context, userID models.UserID) (token *models.Token, fcerr *fcerror.Error) {
	fcerr = authorization.EnforceSelf(authCtx, userID)
	if fcerr != nil {
		return
	}

	token = &models.Token{
		Value:      models.TokenValue(utils.GenerateRandomString(mgr.cfg.GetTokenValueLength())),
		ValidUntil: utils.GetTimeIn(mgr.cfg.GetTokenExpirationDuration()),
		UserID:     userID,
	}

	trans, fcerr := mgr.authPersistence.StartReadWriteTransaction()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create transaction")
		return
	}
	fcerr = trans.SaveToken(token)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to save user")
		trans.Rollback()
		return
	}
	fcerr = trans.Commit()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to commit transaction")
		return
	}
	return token, nil
}
