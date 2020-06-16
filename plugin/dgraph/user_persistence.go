package dgraph

import (
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/config"
	"github.com/freecloudio/server/domain/models"

	"github.com/sirupsen/logrus"
)

func init() {
	persistence.RegisterUserPersistenceController(config.DGraphPersistenceKey, &UserPersistence{})
}

// UserPersistence provides persistence functions for users in dgraph
type UserPersistence struct{}

func (p *UserPersistence) StartTransaction() (persistence.UserPersistenceTransaction, error) {
	return &userTransaction{}, nil
}

type userTransaction struct{}

func (p *userTransaction) SaveUser(user *models.User) (err error) {
	logrus.Print(user)
	return
}

func (p *userTransaction) Commit() (err error) {
	logrus.Print("Commit")
	return
}

func (p *userTransaction) Rollback() (err error) {
	logrus.Print("Rollback")
	return
}
