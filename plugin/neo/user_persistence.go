package neo

import (
	"errors"
	"fmt"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

func init() {
	persistence.RegisterUserPersistenceController(config.NeoPersistenceKey, &UserPersistence{})
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "User", model: &models.User{}})
}

type UserPersistence struct{}

func (up *UserPersistence) StartReadTransaction() (tx persistence.UserPersistenceReadTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeRead)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo read transaction")
		return
	}
	return &userReadTransaction{txCtx}, nil
}

func (up *UserPersistence) StartReadWriteTransaction() (tx persistence.UserPersistenceReadWriteTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeWrite)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo write transaction")
		return
	}
	return &userReadWriteTransaction{userReadTransaction{txCtx}}, nil
}

type userReadTransaction struct {
	*transactionCtx
}

func (tx *userReadTransaction) CountUsers() (count int64, fcerr *fcerror.Error) {
	record, err := neo4j.Single(tx.neoTx.Run(`
		MATCH (u:User)
		RETURN COUNT(u) as count
	`, nil))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrUserNotFound, fcerror.ErrUnknown)
		return
	}

	countInt, ok := record.Get("count")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("count not found in records"))
		return
	}
	count = countInt.(int64)
	return
}

func (tx *userReadTransaction) GetUserByID(userID models.UserID) (user *models.User, fcerr *fcerror.Error) {
	record, err := neo4j.Single(tx.neoTx.Run(`
		MATCH (u:User)
		WHERE ID(u) = $id
		RETURN u
	`, map[string]interface{}{"id": userID}))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrUserNotFound, fcerror.ErrUnknown)
		return
	}

	user = &models.User{}
	fcerr = recordToModel(record, "u", user)
	if fcerr != nil {
		return
	}
	user.ID = userID
	return
}

func (tx *userReadTransaction) GetUserByEmail(email string) (user *models.User, fcerr *fcerror.Error) {
	record, err := neo4j.Single(tx.neoTx.Run(`
		MATCH (u:User)
		WHERE u.email = $email
		RETURN u, ID(u) as id
	`, map[string]interface{}{"email": email}))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrUserNotFound, fcerror.ErrUnknown)
		return
	}

	user = &models.User{}
	fcerr = recordToModel(record, "u", user)
	if fcerr != nil {
		return
	}
	userIDInt, ok := record.Get("id")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("id not found in records"))
		return
	}
	user.ID = models.UserID(userIDInt.(int64))
	return
}

type userReadWriteTransaction struct {
	userReadTransaction
}

func (tx *userReadWriteTransaction) SaveUser(user *models.User) (fcerr *fcerror.Error) {
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
		fcerr = neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
		return
	}

	userIDInt, ok := record.Get("id")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, fmt.Errorf("Failed to convert value to userID: %v", record.GetByIndex(0)))
		return
	}
	user.ID = models.UserID(userIDInt.(int64))

	return
}

func (tx *userReadWriteTransaction) UpdateUser(user *models.User) (fcerr *fcerror.Error) {
	currTime := utils.GetCurrentTime()
	user.Updated = currTime

	result, err := tx.neoTx.Run(`
		MATCH (u:User)
		WHERE ID(u) = $id
		SET u += $user
		`,
		map[string]interface{}{
			"id":   user.ID,
			"user": modelToMap(user),
		})
	if err == nil {
		_, err = result.Consume()
	}
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
		return
	}
	return
}
