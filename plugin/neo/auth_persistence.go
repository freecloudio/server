package neo

import (
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
	persistence.RegisterAuthPersistenceController(config.NeoPersistenceKey, &AuthPersistence{})
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "Session", model: &models.Session{}})
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

func (tx *authReadTransaction) GetSessionByToken(token models.Token) (session *models.Session, fcerr *fcerror.Error) {
	record, err := neo4j.Single(tx.neoTx.Run(`
		MATCH (s:Session)<-[:AUTHENTICATES_WITH]-(u:User)
		WHERE s.token = $token
		RETURN s, ID(u) as user_id
		`,
		map[string]interface{}{
			"token": string(token),
		}))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrTokenNotFound, fcerror.ErrDBReadFailed)
		return
	}

	session = &models.Session{}
	fcerr = recordToModel(record, "s", session)
	if fcerr != nil {
		return
	}

	userIDInt, ok := record.Get("user_id")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, fmt.Errorf("Failed to convert value to userID: %v", record.GetByIndex(0)))
		return
	}
	session.UserID = models.UserID(userIDInt.(int64))

	tx.neoTx.Close()

	return
}

type authReadWriteTransaction struct {
	authReadTransaction
}

func (tx *authReadWriteTransaction) SaveSession(session *models.Session) *fcerror.Error {
	res, err := tx.neoTx.Run(`
		MATCH (u:User)
		WHERE ID(u) = $user_id
		CREATE (u)-[a:AUTHENTICATES_WITH]->(s:Session $s)
		`,
		map[string]interface{}{
			"user_id": session.UserID,
			"s":       modelToMap(session),
		})
	if err == nil {
		_, err = res.Consume()
	}

	return neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
}

func (tx *authReadWriteTransaction) DeleteSessionByToken(token models.Token) *fcerror.Error {
	res, err := tx.neoTx.Run(`
		MATCH (s:Session)
		WHERE s.token = $token
		DETACH DELETE s
		`,
		map[string]interface{}{
			"token": string(token),
		})
	if err == nil {
		_, err = res.Consume()
	}

	return neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
}

func (tx *authReadWriteTransaction) DeleteExpiredSessions() *fcerror.Error {
	res, err := tx.neoTx.Run(`
		MATCH (s:Session)
		WHERE s.valid_until < $now
		DETACH DELETE s
		`,
		map[string]interface{}{
			"now": utils.GetCurrentTime(),
		})
	if err == nil {
		_, err = res.Consume()
	}

	return neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
}
