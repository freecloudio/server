package application

import (
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/sirupsen/logrus"
)

// UserManager contains all use cases related to user management
type UserManager struct{}

func (mgr *UserManager) CreateUser() {
	trans, _ := persistence.GetUserPersistenceController().StartReadWriteTransaction()
	trans.SaveUser(&models.User{FirstName: "Max", LastName: "Heidinger"})
	trans.Commit()
}

func (mgr *UserManager) GetUser(userID models.UserID) {
	trans, _ := persistence.GetUserPersistenceController().StartReadTransaction()
	user, err := trans.GetUser(userID)
	logrus.WithError(err).Printf("Got User: %v", user)
}
