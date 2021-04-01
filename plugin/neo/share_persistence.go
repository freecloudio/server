package neo

import (
	"errors"

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

type shareReadWriteTransaction struct {
	shareReadTransaction
}

func (tx *shareReadWriteTransaction) CreateShare(userID models.UserID, share *models.Share) (created bool, fcerr *fcerror.Error) {
	// userID -> share.nodeID WITHOUT any SHARED
	record, err := neo4j.Single(tx.neoTx.Run(`
		MATCH p = (u:User {id: $user_id})-[:HAS_ROOT_FOLDER|CONTAINS*]->(n:Node {id: $node_id})
		WITH relationships(p)[-1] AS last_relationship
		RETURN last_relationship.name AS node_name
		`,
		map[string]interface{}{
			"user_id": userID,
			"node_id": share.NodeID,
		}))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
		return
	}

	nodeNameInt, ok := record.Get("node_name")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("node_name not found in record"))
		return
	}

	// share.nodeID -> * WITHOUT any SHARED
	res, err := tx.neoTx.Run(`
		MATCH p = (:Node {id: $node_id})-[:CONTAINS|CONTAINS_SHARED*]->(n:Node)
		WITH reduce(x = false, t IN relationships(p) | TYPE(t) = "CONTAINS") AS contains_shared
		WHERE contains_shared = true
		RETURN contains_shared
		`,
		map[string]interface{}{
			"node_id": share.NodeID,
		})
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
		return
	}
	if res.Next() {
		fcerr = fcerror.NewError(fcerror.ErrShareContainsOtherShares, nil)
	}

	// TODO: Check that nodeName is not already used in root folder

	res, err = tx.neoTx.Run(`
			MATCH (u:User {id: $user_id})-[:HAS_ROOT_FOLDER]->(f:Node:Folder)
			MERGE (f)-[r:CONTAINS_SHARED {name: $node_name}]->(n:Node {id: $node_id})
			ON CREATE
				SET r = $share
		`,
		map[string]interface{}{
			"user_id":   share.SharedWithID,
			"node_name": nodeNameInt,
			"node_id":   share.NodeID,
			"share":     modelToMap(share),
		})
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
		return
	}

	summary, err := res.Summary()
	if err == nil && summary.Counters().NodesCreated() > 0 {
		created = true
	}

	fcerr = neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
	return
}
