package application

import (
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
)

// UserManager contains all use cases related to user management
type UserManager struct{}

func (mgr *UserManager) CreateUser() {
	trans, _ := persistence.GetUserPersistenceController().StartReadWriteTransaction()
	trans.SaveUser(&models.User{FirstName: "Max", LastName: "Heidinger"})
	trans.Commit()
}
