package neo

import (
	"errors"
	"fmt"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"
	"github.com/google/uuid"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

type containsRelation struct {
	Name string `json:"name"`
}

func init() {
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "Node", model: &models.Node{}})
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "CONTAINS", model: &containsRelation{}})
}

type NodePersistence struct{}

func CreateNodePersistence(cfg config.Config) (nodePersistence *NodePersistence, fcerr *fcerror.Error) {
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

func (tx *nodeReadTransaction) GetNodeByPath(userID models.UserID, path string, withShared bool) (node *models.Node, fcerr *fcerror.Error) {
	pathSegments := utils.GetPathSegments(path)
	relationCount := len(pathSegments) + 1 // Add one for HAS_ROOT_FOLDER

	relLabels := "HAS_FOLDER|CONTAINS"
	if withShared {
		relLabels += "|CONTAINS_SHARED"
	}

	record, err := neo4j.Single(tx.neoTx.Run(fmt.Sprintf(`
			MATCH p = (u:User {id: $user_id})-[:%s*%d]->(n:Node)
			WHERE [n in tail(relationships(p)) | n.name] = $path_segments
			WITH n, nodes(p)[-2] as second_last_node, relationships(p)[-1] as last_relationship
			RETURN n, "Folder" IN labels(n) AS is_folder, last_relationship.name as name,
				CASE
					WHEN 'Folder' IN labels(second_last_node) THEN second_last_node.id
					ELSE NULL
				END AS parent_node_id
		`, relLabels, relationCount),
		map[string]interface{}{
			"user_id":       userID,
			"path_segments": pathSegments,
		}))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
		return
	}

	return tx.fillNodeInfo(record, userID, path)
}

func (tx *nodeReadTransaction) GetNodeByID(userID models.UserID, nodeID models.NodeID, withShared bool) (node *models.Node, fcerr *fcerror.Error) {
	relLabels := "HAS_FOLDER|CONTAINS"
	if withShared {
		relLabels += "|CONTAINS_SHARED"
	}

	record, err := neo4j.Single(tx.neoTx.Run(fmt.Sprintf(`
			MATCH p = (u:User {id: $user_id})-[:%s*]->(n:Node {id: $node_id})
			WITH n, p, nodes(p)[-2] as second_last_node, relationships(p)[-1] as last_relationship
			RETURN n, "Folder" IN labels(n) AS is_folder,
				reduce(s = "", n in tail(relationships(p)) | s + '/' + n.name) as path,
				last_relationship.name as name,
				CASE
					WHEN 'Folder' IN labels(second_last_node) THEN second_last_node.id
					ELSE NULL
				END AS parent_node_id
		`, relLabels),
		map[string]interface{}{
			"user_id": userID,
			"node_id": nodeID,
		}))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
		return
	}

	pathInt, ok := record.Get("path")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("path not found in record"))
		return
	}
	path := pathInt.(string)

	return tx.fillNodeInfo(record, userID, path)
}

func (tx *nodeReadTransaction) fillNodeInfo(record neo4j.Record, userID models.UserID, path string) (node *models.Node, fcerr *fcerror.Error) {
	node = &models.Node{}
	fcerr = recordToModel(record, "n", node)
	if fcerr != nil {
		return
	}

	if path == "" {
		path = "/"
	}
	node.Path, _ = utils.SplitPath(path)
	node.FullPath = path
	node.PerspectiveUserID = userID

	nameInt, ok := record.Get("name")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("name not found in record"))
		return
	}
	node.Name, _ = nameInt.(string)

	isFolderInt, ok := record.Get("is_folder")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("is_folder not found in record"))
		return
	}
	if isFolder := isFolderInt.(bool); isFolder {
		node.Type = models.NodeTypeFolder
	} else {
		node.Type = models.NodeTypeFile
	}

	parentNodeIDInt, ok := record.Get("parent_node_id")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("parent_node_id not found in record"))
		return
	}
	if parentNodeIDStr, ok := parentNodeIDInt.(string); ok {
		parentNodeID := models.NodeID(parentNodeIDStr)
		node.ParentNodeID = &parentNodeID
	}

	node.OwnerID, fcerr = tx.getOwnerOfNodeID(node.ID)
	if fcerr != nil {
		return
	}

	// TODO: Insert ShareMode, Starred
	return
}

func (tx *nodeReadTransaction) getOwnerOfNodeID(nodeID models.NodeID) (userID models.UserID, fcerr *fcerror.Error) {
	record, err := neo4j.Single(tx.neoTx.Run(`
			MATCH (u:User)-[:HAS_ROOT_FOLDER|CONTAINS*]->(n:Node {id: $node_id})
			RETURN u.id as id
		`,
		map[string]interface{}{
			"node_id": nodeID,
		}))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
		return
	}

	userIDInt, ok := record.Get("id")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("id not found in record"))
		return
	}
	userID = models.UserID(userIDInt.(string))
	return
}

type nodeReadWriteTransaction struct {
	nodeReadTransaction
}

func (tx *nodeReadWriteTransaction) CreateUserRootFolder(userID models.UserID) (created bool, fcerr *fcerror.Error) {
	insertNode := &models.Node{
		ID:      models.NodeID(uuid.NewString()),
		Created: utils.GetCurrentTime(),
		Updated: utils.GetCurrentTime(),
	}

	res, err := tx.neoTx.Run(`
		MATCH (u:User {id: $user_id})
		MERGE (u)-[:HAS_ROOT_FOLDER]->(f:Node:Folder)
		ON CREATE
			SET f = $f
		`,
		map[string]interface{}{
			"user_id": userID,
			"f":       modelToMap(insertNode),
		})
	if err == nil {
		var summary neo4j.ResultSummary
		summary, err = res.Consume()
		if summary.Counters().NodesCreated() > 0 {
			created = true
		}
	}

	fcerr = neoToFcError(err, fcerror.ErrUnknown, fcerror.ErrDBWriteFailed)
	return
}

func (tx *nodeReadWriteTransaction) CreateNodeByID(userID models.UserID, nodeType models.NodeType, parentNodeID models.NodeID, name string) (node *models.Node, created bool, fcerr *fcerror.Error) {
	insertNode := &models.Node{
		ID:      models.NodeID(uuid.NewString()),
		Created: utils.GetCurrentTime(),
		Updated: utils.GetCurrentTime(),
	}
	insertRelation := &containsRelation{
		Name: name,
	}
	insertNodeType := "File"
	if nodeType == models.NodeTypeFolder {
		insertNodeType = "Folder"
	}

	result, err := tx.neoTx.Run(fmt.Sprintf(`
			MATCH p = (u:User {id: $user_id})-[:HAS_ROOT_FOLDER|CONTAINS|CONTAINS_SHARED*]->(f:Node:Folder {id: $parent_node_id})
			MERGE (f)-[r:CONTAINS {name: $r.name}]->(n:Node)
			ON CREATE
				SET n:%s
				SET n = $n
				SET r = $r
			WITH n, r, p
			RETURN n,
				"Folder" IN labels(n) AS is_folder,
				reduce(s = "", n in tail(relationships(p)) | s + '/' + n.name) as parent_path,
				r.name as name,
				$parent_node_id AS parent_node_id
		`, insertNodeType),
		map[string]interface{}{
			"user_id":        userID,
			"parent_node_id": parentNodeID,
			"n":              modelToMap(insertNode),
			"r":              modelToMap(insertRelation),
		})
	record, err := neo4j.Single(result, err)
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
		return
	}

	summary, err := result.Summary()
	if err == nil && summary.Counters().NodesCreated() > 0 {
		created = true
	}

	parentPathInt, ok := record.Get("parent_path")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("parent_path not found in record"))
		return
	}
	path := parentPathInt.(string)
	if path == "" {
		path = "/"
	}
	path = utils.JoinPaths(path, name)

	node, fcerr = tx.fillNodeInfo(record, userID, path)
	return
}
