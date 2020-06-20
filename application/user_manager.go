package application

import (
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/sirupsen/logrus"
)

// UserManager contains all use cases related to user management
type UserManager struct{}

func (mgr *UserManager) CreateUser() (err error) {
	trans, err := persistence.GetUserPersistenceController().StartReadWriteTransaction()
	if err != nil {
		logrus.WithError(err).Error("Failed to create transaction")
		return
	}
	err = trans.SaveUser(&models.User{FirstName: "Max", LastName: "Heidinger"})
	if err != nil {
		logrus.WithError(err).Error("Failed to save user")
		return
	}
	err = trans.Commit()
	if err != nil {
		logrus.WithError(err).Error("Failed to commit transaction")
		return
	}
	return
}

func (mgr *UserManager) GetUser(userID models.UserID) (user *models.User, err error) {
	trans, err := persistence.GetUserPersistenceController().StartReadTransaction()
	if err != nil {
		logrus.WithError(err).Error("Failed to create transaction")
		return
	}
	user, err = trans.GetUser(userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user")
		return
	}
	return
}
