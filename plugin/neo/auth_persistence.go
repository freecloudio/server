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
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "Session", model: &models.Session{}})
}

type AuthPersistence struct{}

func CreateAuthPersistence(cfg config.Config) (authPersistence *AuthPersistence, fcerr *fcerror.Error) {
	if neo == nil {
		fcerr = initializeNeo(cfg)
		if fcerr != nil {
			return
		}
	}
	authPersistence = &AuthPersistence{}
	return
}

func (*AuthPersistence) Close() *fcerror.Error {
	if neo != nil {
		return closeNeo()
	}
	return nil
}

func (*AuthPersistence) StartReadTransaction() (tx persistence.AuthPersistenceReadTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeRead)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo read transaction")
		return
	}
	return &authReadTransaction{txCtx}, nil
}

func (*AuthPersistence) StartReadWriteTransaction() (tx persistence.AuthPersistenceReadWriteTransaction, fcerr *fcerror.Error) {
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
		MATCH (s:Session {token: $token})<-[:AUTHENTICATES_WITH]-(u:User)
		RETURN s, u.id as user_id
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
	session.UserID = models.UserID(userIDInt.(string))

	return
}

type authReadWriteTransaction struct {
	authReadTransaction
}

func (tx *authReadWriteTransaction) SaveSession(session *models.Session) *fcerror.Error {
	res, err := tx.neoTx.Run(`
		MATCH (u:User {id: $user_id})
		CREATE (u)-[:AUTHENTICATES_WITH]->(:Session $s)
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
		MATCH (s:Session {token: $token})
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
