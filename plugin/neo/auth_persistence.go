package neo

import (
	"fmt"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

func init() {
	persistence.RegisterAuthPersistenceController(config.NeoPersistenceKey, &AuthPersistence{})
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "Token", model: &models.Token{}})
}

type AuthPersistence struct{}

func (up *AuthPersistence) StartReadTransaction() (tx persistence.AuthPersistenceReadTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeRead)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo read transaction")
		return
	}
	return &authReadTransaction{txCtx}, nil
}

func (up *AuthPersistence) StartReadWriteTransaction() (tx persistence.AuthPersistenceReadWriteTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeWrite)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo write transaction")
		return
	}
	return &authReadWriteTransaction{authReadTransaction{txCtx}}, nil
}

type authReadTransaction struct {
	*transactionCtx
}

func (tx *authReadTransaction) CheckToken(tokenValue models.TokenValue) (token *models.Token, fcerr *fcerror.Error) {
	record, err := neo4j.Single(tx.neoTx.Run(`
		MATCH (t:Token)<-[:AUTHENTICATES_WITH]-(u:User)
		WHERE t.value = $token_value
		RETURN t, ID(u) as user_id
		`,
		map[string]interface{}{
			"token_value": string(tokenValue),
		}))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrTokenNotFound, fcerror.ErrDBReadFailed)
		return
	}

	token = &models.Token{}
	err = recordToModel(record, "t", token)
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, err)
		return
	}

	userIDInt, ok := record.Get("user_id")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, fmt.Errorf("Failed to convert value to userID: %v", record.GetByIndex(0)))
		return
	}
	token.UserID = models.UserID(userIDInt.(int64))

	tx.neoTx.Close()

	return
}

type authReadWriteTransaction struct {
	authReadTransaction
}

func (tx *authReadWriteTransaction) SaveToken(token *models.Token) *fcerror.Error {
	_, err := tx.neoTx.Run(`
		MATCH (u:User)
		WHERE ID(u) = $t.user_id
		CREATE (u)-[a:AUTHENTICATES_WITH]->(t:Token { value: $t.value, valid_until: $t.valid_until })
		`,
		map[string]interface{}{
			"t": modelToMap(token),
		})

	return neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
}

func (tx *authReadWriteTransaction) DeleteToken(tokenValue models.TokenValue) *fcerror.Error {
	_, err := tx.neoTx.Run(`
		MATCH (t:Token)
		WHERE t.value = $token_value
		DETACH DELETE t
		`,
		map[string]interface{}{
			"token_value": string(tokenValue),
		})

	return neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
}
