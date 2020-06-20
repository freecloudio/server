package neo

import (
	"fmt"

	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/config"
	"github.com/freecloudio/server/domain/models"
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
	}
	return &userReadTransaction{txCtx}, nil
}

func (up *UserPersistence) StartReadWriteTransaction() (persistence.UserPersistenceReadWriteTransaction, error) {
	txCtx, err := newTransactionContext(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("Failed to create neo write transaction")
	}
	return &userReadWriteTransaction{txCtx}, nil
}

type userReadTransaction struct {
	*transactionCtx
}

func (tx *userReadTransaction) GetUser(userID models.UserID) (*models.User, error) {
	return nil, nil
}

type userReadWriteTransaction struct {
	*transactionCtx
}

func (tx *userReadWriteTransaction) SaveUser(user *models.User) (err error) {
	res, err := tx.neoTx.Run(`
		CREATE (u:User { 
			first_name: $first_name,
			last_name: $last_name,
			email: $email,
			password: $password,
			is_admin: $is_admin,
			created: $created,
			updated: $updated
		})
		RETURN id(u) as id
		`, map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"password":   user.Password,
		"is_admin":   user.IsAdmin,
		"created":    user.Created,
		"updated":    user.Updated,
	})
	if err != nil {
		return err
	}

	if res.Next() {
		userID, ok := res.Record().GetByIndex(0).(models.UserID)
		if !ok {
			err = fmt.Errorf("Failed to convert value to userID: %v", res.Record().GetByIndex(0))
			return
		}
		user.ID = userID
	} else {
		err = fmt.Errorf("Result does not contain result row")
		return
	}

	return
}

func (tx *userReadWriteTransaction) Commit() (err error) {
	txErr := tx.neoTx.Commit()
	if txErr != nil {
		logrus.WithError(txErr).Error("Failed to commit neo transaction - close session anyway")
		err = txErr
	}
	sessErr := tx.session.Close()
	if sessErr != nil {
		logrus.WithError(sessErr).Error("Failed to close neo session")
		if err == nil {
			err = sessErr
		}
	}
	return
}

func (tx *userReadWriteTransaction) Rollback() (err error) {
	txErr := tx.neoTx.Rollback()
	if txErr != nil {
		logrus.WithError(txErr).Error("Failed to rollback neo transaction - close session anyway")
		err = txErr
	}
	sessErr := tx.session.Close()
	if sessErr != nil {
		logrus.WithError(sessErr).Error("Failed to close neo session")
		if err == nil {
			err = sessErr
		}
	}
	return
}
