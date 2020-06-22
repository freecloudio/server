package application

import (
	"time"

	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/sirupsen/logrus"
)

type AuthManager struct{}

func (mgr *AuthManager) CreateToken(userID models.UserID) (token *models.Token, err error) {
	token = &models.Token{
		Value:      models.TokenValue("AAAABBBB"),
		ValidUntil: time.Now().Add(time.Hour).UTC(),
		UserID:     userID,
	}

	trans, err := persistence.GetAuthPersistenceController().StartReadWriteTransaction()
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
