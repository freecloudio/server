package neo

import (
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

type SharePersistence struct{}

func init() {
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "CONTAINS_SHARED", model: &containsRelation{}})
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "CONTAINS_SHARED", model: &models.Share{}})
}

func CreateSharePersistence(cfg config.Config) (sharePersistence *SharePersistence, fcerr *fcerror.Error) {
	if neo == nil {
		fcerr = initializeNeo(cfg)
		if fcerr != nil {
			return
		}
	}
	sharePersistence = &SharePersistence{}
	return
}

func (*SharePersistence) Close() *fcerror.Error {
	if neo != nil {
		return closeNeo()
	}
	return nil
}

func (*SharePersistence) StartReadTransaction() (tx persistence.SharePersistenceReadTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeRead)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo read transaction")
		return
	}
	return &shareReadTransaction{txCtx}, nil
}

func (*SharePersistence) StartReadWriteTransaction() (tx persistence.SharePersistenceReadWriteTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeWrite)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo write transaction")
		return
	}
	return &shareReadWriteTransaction{shareReadTransaction{txCtx}}, nil
}

type shareReadTransaction struct {
	*transactionCtx
}

func (tx *shareReadTransaction) NodeContainsNestedShares(nodeID models.NodeID) (containsShared bool, fcerr *fcerror.Error) {
	res, err := tx.neoTx.Run(`
		MATCH p = (:Node {id: $node_id})-[:CONTAINS|CONTAINS_SHARED*]->(n:Node)
		WITH reduce(x = false, t IN relationships(p) | TYPE(t) = "CONTAINS_SHARED") AS contains_shared
		WHERE contains_shared = true
		RETURN contains_shared
		`,
		map[string]interface{}{
			"node_id": nodeID,
		})
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
		return
	}
	if res.Next() {
		containsShared = true
	}
	return
}

type shareReadWriteTransaction struct {
	shareReadTransaction
}

func (tx *shareReadWriteTransaction) CreateShare(userID models.UserID, share *models.Share, insertName string) (created bool, fcerr *fcerror.Error) {
	res, err := tx.neoTx.Run(`
			MATCH (u:User {id: $user_id})-[:HAS_ROOT_FOLDER]->(f:Node:Folder), (n:Node {id: $node_id})
			MERGE (f)-[r:CONTAINS_SHARED {name: $node_name}]->(n)
			ON CREATE
				SET r += $share
		`,
		map[string]interface{}{
			"user_id":   share.SharedWithID,
			"node_name": insertName,
			"node_id":   share.NodeID,
			"share":     modelToMap(share),
		})
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
		return
	}

	summary, err := res.Summary()
	if err == nil && summary.Counters().RelationshipsCreated() > 0 {
		created = true
	}

	fcerr = neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
	return
}
