package gin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	pathParam          = "path"
	nodeIDParam        = "node_id"
	filenameParam      = "filename"
	nodeTypeQueryParam = "node_type"
)

func (r *Router) buildNodeRoutes() {
	grp := r.engine.Group("/api/node")

	grp.GET("info/path/*"+pathParam, r.getNodeInfoByPath)
	grp.GET("info/id/:"+nodeIDParam, r.getNodeInfoByID)
	grp.POST(fmt.Sprintf("create/id/:%s/:%s", nodeIDParam, filenameParam), r.createNodeByID)
	// list/id/
	// content/id/
}

func (r *Router) getNodeInfoByPath(c *gin.Context) {
	authContext := getAuthContext(c)
	path := c.Param(pathParam)

	node, fcerr := r.managers.Node.GetNodeByPath(authContext, path)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (r *Router) getNodeInfoByID(c *gin.Context) {
	authContext := getAuthContext(c)
	nodeID, fcerr := extractNodeID(c)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get nodeID from request")
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	node, fcerr := r.managers.Node.GetNodeByID(authContext, nodeID)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (r *Router) createNodeByID(c *gin.Context) {
	authContext := getAuthContext(c)
	nodeID, fcerr := extractNodeID(c)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get nodeID from request")
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}
	filename := c.Param(filenameParam)

	nodeTypeStr := c.Query(nodeTypeQueryParam)
	nodeType := models.NodeTypeFolder
	if nodeTypeStr == string(models.NodeTypeFile) {
		nodeType = models.NodeTypeFile
	}

	createdNodeID, fcerr := r.managers.Node.CreateNode(authContext, nodeType, nodeID, filename)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"created_node_id": createdNodeID})
}

func extractNodeID(c *gin.Context) (nodeID models.NodeID, fcerr *fcerror.Error) {
	nodeIDStr := c.Param(nodeIDParam)
	if nodeIDStr == "" {
		fcerr = fcerror.NewError(fcerror.ErrBadRequest, errors.New("NodeID not found in path param"))
		return
	}

	nodeIDInt, err := strconv.Atoi(nodeIDStr)
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrBadRequest, err)
		return
	}

	nodeID = models.NodeID(nodeIDInt)
	return
}
