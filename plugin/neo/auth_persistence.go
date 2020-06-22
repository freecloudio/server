package neo

import (
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/config"
	"github.com/freecloudio/server/domain/models"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

func init() {
	persistence.RegisterAuthPersistenceController(config.NeoPersistenceKey, &AuthPersistence{})
}

type AuthPersistence struct{}

func (up *AuthPersistence) StartReadTransaction() (tx persistence.AuthPersistenceReadTransaction, err error) {
	txCtx, err := newTransactionContext(neo4j.AccessModeRead)
	if err != nil {
		logrus.WithError(err).Error("Failed to create neo read transaction")
		return
	}
	return &authReadTransaction{txCtx}, nil
}

func (up *AuthPersistence) StartReadWriteTransaction() (tx persistence.AuthPersistenceReadWriteTransaction, err error) {
	txCtx, err := newTransactionContext(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("Failed to create neo write transaction")
		return
	}
	return &authReadWriteTransaction{txCtx}, nil
}

type authReadTransaction struct {
	*transactionCtx
}

func (tx *authReadTransaction) CheckToken(tokenValue models.TokenValue) (token *models.Token, err error) {
	return
}

type authReadWriteTransaction struct {
	*transactionCtx
}

func (tx *authReadWriteTransaction) SaveToken(token *models.Token) (err error) {
	_, err = tx.neoTx.Run(`
		MATCH (u:User)
		WHERE ID(u) = $t.user_id
		CREATE (u)-[a:AUTHENTICATES_WITH]->(t:Token { value: $t.value, valid_until: $t.valid_until })
		`,
		map[string]interface{}{
			"t": modelToMap(token),
		})

	return
}
