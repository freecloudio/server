package neo

import (
	"fmt"

	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/config"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func init() {
	persistence.RegisterPluginInitialization(config.NeoPersistenceKey, InitializeNeo)
}

var neo neo4j.Driver

func InitializeNeo() (err error) {
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "freecloud", ""), noEncrypted)
	if err != nil {
		return
	}
	neo = driver
	return
}

func noEncrypted(config *neo4j.Config) {
	config.Encrypted = false
}

type transactionCtx struct {
	session neo4j.Session
	neoTx   neo4j.Transaction
}

func newTransactionContext(accessMode neo4j.AccessMode) (txCtx *transactionCtx, err error) {
	session, err := neo.NewSession(neo4j.SessionConfig{AccessMode: accessMode})
	if err != nil {
		err = fmt.Errorf("failed to create neo4j session: %v", err)
		return
	}
	neoTx, err := session.BeginTransaction()
	if err != nil {
		err = fmt.Errorf("failed to create neo4j transaction: %v", err)
		session.Close()
		return
	}
	txCtx = &transactionCtx{session, neoTx}
	return
}
