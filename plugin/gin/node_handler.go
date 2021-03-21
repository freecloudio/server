package gin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	pathParam      = "path"
	nodeIDParam    = "node_id"
	filenameParam  = "filename"
	fileQueryParam = "file"
)

func (r *Router) buildNodeRoutes() {
	grp := r.engine.Group("/api/node")

	grp.GET("info/path/*"+pathParam, r.getNodeInfoByPath)
	grp.GET("info/id/:"+nodeIDParam, r.getNodeInfoByID)
	grp.POST(fmt.Sprintf("create/id/:%s/:%s", nodeIDParam, filenameParam), r.createNodeByID)
	grp.POST("/upload/id/:"+nodeIDParam, r.uploadFileByID)
	grp.POST(fmt.Sprintf("upload/id/:%s/:%s", nodeIDParam, filenameParam), r.uploadFileByParentID)
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

	_, file := c.GetQuery(fileQueryParam)
	nodeType := models.NodeTypeFolder
	if file {
		nodeType = models.NodeTypeFile
	}

	createdNode, created, fcerr := r.managers.Node.CreateNode(authContext, nodeType, nodeID, filename)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"node_id": createdNode.ID, "created": created})
}

func (r *Router) uploadFileByParentID(c *gin.Context) {
	authContext := getAuthContext(c)
	parentNodeID, fcerr := extractNodeID(c)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get parentNodeID from request")
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logrus.WithError(err).Error("No file attached to upload")
		fcerr = fcerror.NewError(fcerror.ErrBadRequest, err)
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}
	filename := c.Param(filenameParam)

	tmpPath := utils.JoinPaths(r.cfg.GetFileStorageTempBasePath(), utils.GenerateRandomString(10))
	err = c.SaveUploadedFile(file, tmpPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to save upload to temp file")
		fcerr = fcerror.NewError(fcerror.ErrCopyFileFailed, err)
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	createdNode, created, fcerr := r.managers.Node.UploadFile(authContext, parentNodeID, filename, tmpPath)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"node_id": createdNode.ID, "created": created})
}

func (r *Router) uploadFileByID(c *gin.Context) {
	authContext := getAuthContext(c)
	nodeID, fcerr := extractNodeID(c)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get nodeID from request")
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logrus.WithError(err).Error("No file attached to upload")
		fcerr = fcerror.NewError(fcerror.ErrBadRequest, err)
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	tmpPath := utils.JoinPaths(r.cfg.GetFileStorageTempBasePath(), utils.GenerateRandomString(10))
	err = c.SaveUploadedFile(file, tmpPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to save upload to temp file")
		fcerr = fcerror.NewError(fcerror.ErrCopyFileFailed, err)
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	fcerr = r.managers.Node.UploadFileByID(authContext, nodeID, tmpPath)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, &gin.H{})
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
