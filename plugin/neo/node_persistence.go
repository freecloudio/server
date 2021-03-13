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
	nodeModelMappings = append(nodeModelMappings, &labelModelMapping{label: "Node", model: &models.Node{}})
}

type NodePersistence struct{}

func CreateNodeePersistence(cfg config.Config) (nodePersistence *NodePersistence, fcerr *fcerror.Error) {
	if neo == nil {
		fcerr = initializeNeo(cfg)
		if fcerr != nil {
			return
		}
	}
	nodePersistence = &NodePersistence{}
	return
}

func (*NodePersistence) Close() *fcerror.Error {
	if neo != nil {
		return closeNeo()
	}
	return nil
}

func (*NodePersistence) StartReadTransaction() (tx persistence.NodePersistenceReadTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeRead)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo read transaction")
		return
	}
	return &nodeReadTransaction{txCtx}, nil
}

func (*NodePersistence) StartReadWriteTransaction() (tx persistence.NodePersistenceReadWriteTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeWrite)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo write transaction")
		return
	}
	return &nodeReadWriteTransaction{nodeReadTransaction{txCtx}}, nil
}

type nodeReadTransaction struct {
	*transactionCtx
}

func (tx *nodeReadTransaction) GetNodeByPath(userID models.UserID, path string) (node *models.Node, fcerr *fcerror.Error) {
	pathSegments := utils.GetPathSegments(path)

	record, err := neo4j.Single(tx.neoTx.Run(fmt.Sprintf(`
			MATCH p = (u:User)-[:HAS_ROOT_FOLDER|CONTAINS|CONTAINS_SHARED*%d]->(n:Node)
			WHERE ID(u) = $user_id AND [n in tail(tail(nodes(p))) | n.name] = $path_segments
			RETURN n, "Folder" IN labels(n) AS is_folder
		`, len(pathSegments)),
		map[string]interface{}{
			"user_id":       userID,
			"path_segments": pathSegments,
		}))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
		return
	}

	node = &models.Node{}
	fcerr = recordToModel(record, "n", node)
	if fcerr != nil {
		return
	}

	isFolderInt, ok := record.Get("is_folder")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("Failed to get is_folder of session"))
		return
	}
	if isFolder := isFolderInt.(bool); isFolder {
		node.Type = models.NodeTypeFolder
	} else {
		node.Type = models.NodeTypeFile
	}

	node.Path, _ = utils.SplitPath(path)
	node.FullPath = path
	node.PerspectiveUserID = userID

	// TODO: Insert ShareMode, Starred
	return
}

type nodeReadWriteTransaction struct {
	nodeReadTransaction
}

func (tx *nodeReadWriteTransaction) CreateUserRootFolder(userID models.UserID) (fcerr *fcerror.Error) {
	insertNode := &models.Node{
		Created: utils.GetCurrentTime(),
		Updated: utils.GetCurrentTime(),
		OwnerID: userID,
		Name:    "",
	}

	res, err := tx.neoTx.Run(`
		MATCH (u:User)
		WHERE ID(u) = $user_id
		MERGE (u)-[:HAS_ROOT_FOLDER]->(f:Node:Folder)
		ON CREATE
			SET f = $f
		`,
		map[string]interface{}{
			"user_id": userID,
			"f":       modelToMap(insertNode),
		})
	if err == nil {
		_, err = res.Consume()
	}

	return neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
}
