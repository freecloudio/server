package neo

import (
	"fmt"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/utils"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

func init() {
	persistence.RegisterUserPersistenceController(config.NeoPersistenceKey, &UserPersistence{})
}

type UserPersistence struct{}

func (up *UserPersistence) StartReadTransaction() (tx persistence.UserPersistenceReadTransaction, err error) {
	txCtx, err := newTransactionContext(neo4j.AccessModeRead)
	if err != nil {
		logrus.WithError(err).Error("Failed to create neo read transaction")
		return
	}
	return &userReadTransaction{txCtx}, nil
}

func (up *UserPersistence) StartReadWriteTransaction() (tx persistence.UserPersistenceReadWriteTransaction, err error) {
	txCtx, err := newTransactionContext(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("Failed to create neo write transaction")
		return
	}
	return &userReadWriteTransaction{txCtx}, nil
}

type userReadTransaction struct {
	*transactionCtx
}

func (tx *userReadTransaction) GetUser(userID models.UserID) (user *models.User, err error) {
	record, err := neo4j.Single(tx.neoTx.Run(`
		MATCH (u:User)
		WHERE ID(u) = $id
		RETURN u
	`, map[string]interface{}{"id": userID}))
	if err != nil {
		return
	}

	user = &models.User{}
	err = recordToModel(record, "u", user)
	return
}

type userReadWriteTransaction struct {
	*transactionCtx
}

func (tx *userReadWriteTransaction) SaveUser(user *models.User) (err error) {
	currTime := utils.GetCurrentTime()
	user.Created = currTime
	user.Updated = currTime

	record, err := neo4j.Single(tx.neoTx.Run(`
		CREATE (u:User $user)
		RETURN ID(u) as id
		`,
		map[string]interface{}{
			"user": modelToMap(user),
		}))
	if err != nil {
		return err
	}

	userIDInt, ok := record.GetByIndex(0).(int64)
	if !ok {
		err = fmt.Errorf("Failed to convert value to userID: %v", record.GetByIndex(0))
		return
	}
	user.ID = models.UserID(userIDInt)

	return
}
