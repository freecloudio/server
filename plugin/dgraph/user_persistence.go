package dgraph

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/config"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/plugin/dgraph/schema"
	"github.com/sirupsen/logrus"
)

func init() {
	persistence.RegisterUserPersistenceController(config.DGraphPersistenceKey, &UserPersistence{})
}

// UserPersistence provides persistence functions for users in dgraph
type UserPersistence struct{}

func (p *UserPersistence) StartReadTransaction() (persistence.UserPersistenceReadTransaction, error) {
	return &userReadTransaction{
		dgraphTx: dg.NewReadOnlyTxn(),
	}, nil
}

func (p *UserPersistence) StartReadWriteTransaction() (persistence.UserPersistenceReadWriteTransaction, error) {
	return &userReadWriteTransaction{
		dgraphTx: dg.NewTxn(),
	}, nil
}

type userReadTransaction struct {
	dgraphTx *dgo.Txn
}

func (tx *userReadTransaction) GetUser(userID models.UserID) (user *models.User, err error) {
	return
}

type userReadWriteTransaction struct {
	dgraphTx *dgo.Txn
}

func (tx *userReadWriteTransaction) SaveUser(user *models.User) (err error) {
	insertUser := schema.CreateDUser(user)
	insertJSON, err := json.Marshal(insertUser)
	if err != nil {
		logrus.WithError(err).WithField("user", user).Error("Failed to SaveUser")
		return errors.New("failed to SaveUser in DGraph")
	}

	response, err := tx.dgraphTx.Mutate(context.TODO(), &api.Mutation{SetJson: insertJSON})

	logrus.WithField("response", response).WithError(err).Print("SaveUser Mutation Done")

	return
}

func (tx *userReadWriteTransaction) Commit() (err error) {
	return tx.dgraphTx.Commit(context.TODO())
}

func (tx *userReadWriteTransaction) Rollback() (err error) {
	return tx.dgraphTx.Discard(context.TODO())
}
