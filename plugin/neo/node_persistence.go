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
)

// TODO: Deduplicate Cypher statements

type containsRelation struct {
	Name string `json:"name"`
}

func getContainsRelationshipLabels(includedShareMode models.ShareMode) string {
	// TODO: Read, Read&Write
	relLabels := "HAS_ROOT_FOLDER|CONTAINS"
	if includedShareMode != models.ShareModeNone {
		relLabels += "|CONTAINS_SHARED"
	}
	return relLabels
}

func init() {
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "Node", model: &models.Node{}})
	labelModelMappings = append(labelModelMappings, &labelModelMapping{label: "CONTAINS", model: &containsRelation{}})
}

type NodePersistence struct {
	logger utils.Logger
}

func CreateNodePersistence(cfg config.Config) (nodePersistence *NodePersistence, fcerr *fcerror.Error) {
	if neo == nil {
		fcerr = initializeNeo(cfg)
		if fcerr != nil {
			return
		}
	}
	nodePersistence = &NodePersistence{logger: utils.CreateLogger(cfg.GetLoggingConfig())}
	return
}

func (*NodePersistence) Close() *fcerror.Error {
	if neo != nil {
		return closeNeo()
	}
	return nil
}

func (p *NodePersistence) StartReadTransaction() (tx persistence.NodePersistenceReadTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeRead, p.logger)
	if fcerr != nil {
		p.logger.WithError(fcerr).Error("Failed to create neo read transaction")
		return
	}
	return &nodeReadTransaction{txCtx}, nil
}

func (p *NodePersistence) StartReadWriteTransaction() (tx persistence.NodePersistenceReadWriteTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeWrite, p.logger)
	if fcerr != nil {
		p.logger.WithError(fcerr).Error("Failed to create neo write transaction")
		return
	}
	return &nodeReadWriteTransaction{nodeReadTransaction{txCtx}}, nil
}

type nodeReadTransaction struct {
	*transactionCtx
}

func (tx *nodeReadTransaction) GetNodeByPath(userID models.UserID, path string, includedShareMode models.ShareMode) (node *models.Node, fcerr *fcerror.Error) {
	pathSegments := utils.GetPathSegments(path)
	relationCount := len(pathSegments) + 1 // Add one for HAS_ROOT_FOLDER
	relLabels := getContainsRelationshipLabels(includedShareMode)

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

	node = &models.Node{}
	fcerr = tx.fillNodeInfo(node, record, userID, path)
	return node, fcerr
}

func (tx *nodeReadTransaction) GetNodeByID(userID models.UserID, nodeID models.NodeID, includedShareMode models.ShareMode) (node *models.Node, fcerr *fcerror.Error) {
	relLabels := getContainsRelationshipLabels(includedShareMode)

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

	node = &models.Node{}
	fcerr = tx.fillNodeInfo(node, record, userID, path)
	return node, fcerr
}

func (tx *nodeReadTransaction) ListByID(userID models.UserID, nodeID models.NodeID, includedShareMode models.ShareMode) (list []*models.Node, fcerr *fcerror.Error) {
	relLabels := getContainsRelationshipLabels(includedShareMode)

	res, err := tx.neoTx.Run(fmt.Sprintf(`
			MATCH p = (u:User {id: $user_id})-[:%s*]->(:Node:Folder {id: $node_id})-[:%s]->(n:Node)
			WITH n, p, nodes(p)[-2] as second_last_node, relationships(p)[-1] as last_relationship
			RETURN n, "Folder" IN labels(n) AS is_folder,
				reduce(s = "", n in tail(relationships(p)) | s + '/' + n.name) as path,
				last_relationship.name as name,
				CASE
					WHEN 'Folder' IN labels(second_last_node) THEN second_last_node.id
					ELSE NULL
				END AS parent_node_id
		`, relLabels, relLabels),
		map[string]interface{}{
			"user_id": userID,
			"node_id": nodeID,
		})
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
		return
	}

	for res.Next() {
		record := res.Record()

		pathInt, ok := record.Get("path")
		if !ok {
			fcerr = fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("path not found in record"))
			return
		}
		path := pathInt.(string)

		node := &models.Node{}
		fcerr := tx.fillNodeInfo(node, record, userID, path)
		if fcerr != nil {
			fcerr = neoToFcError(err, fcerror.ErrNodeNotFound, fcerror.ErrDBReadFailed)
			return nil, fcerr
		}

		list = append(list, node)
	}
	return
}

func (tx *nodeReadTransaction) fillNodeInfo(node *models.Node, record neo4j.Record, userID models.UserID, path string) (fcerr *fcerror.Error) {
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

	// TODO: Insert Starred
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

func (tx *nodeReadWriteTransaction) CreateNodeByID(userID models.UserID, node *models.Node) (created bool, fcerr *fcerror.Error) {
	node.ID = models.NodeID(uuid.NewString())
	node.Created = utils.GetCurrentTime()
	node.Updated = utils.GetCurrentTime()
	insertRelation := &containsRelation{
		Name: node.Name,
	}
	insertNodeType := "File"
	if node.Type == models.NodeTypeFolder {
		insertNodeType = "Folder"
	}

	result, err := tx.neoTx.Run(fmt.Sprintf(`
			MATCH p = (u:User {id: $user_id})-[:HAS_ROOT_FOLDER|CONTAINS|CONTAINS_SHARED*]->(f:Node:Folder {id: $parent_node_id})
			MERGE (f)-[r:CONTAINS {name: $r.name}]->(n:Node)
			ON CREATE
				SET n:%s
				SET n += $n
				SET r += $r
			WITH n, r, p
			RETURN n,
				"Folder" IN labels(n) AS is_folder,
				reduce(s = "", n in tail(relationships(p)) | s + '/' + n.name) as parent_path,
				r.name as name,
				$parent_node_id AS parent_node_id
		`, insertNodeType),
		map[string]interface{}{
			"user_id":        userID,
			"parent_node_id": node.ParentNodeID,
			"n":              modelToMap(node),
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
	path = utils.JoinPaths(path, node.Name)

	fcerr = tx.fillNodeInfo(node, record, userID, path)
	return
}
