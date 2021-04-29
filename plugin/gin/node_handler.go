package gin

import (
	"errors"
	"net/http"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	pathParam     = "path"
	nodeIDParam   = "node_id"
	fileNameParam = "filename"
)

func (r *Router) buildNodeRoutes() {
	grp := r.engine.Group("/api/node")

	grp.GET("info/path/*"+pathParam, r.getNodeInfoByPath)
	grp.GET("info/id/:"+nodeIDParam, r.getNodeInfoByID)
	grp.GET("list/id/:"+nodeIDParam, r.getNodeListByID)
	grp.GET("content/id/:"+nodeIDParam, r.getNodeContentByID)
	grp.POST("create/", r.createNodeByID)
	grp.POST("upload/id/:"+nodeIDParam, r.uploadFileByID)
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

func (r *Router) getNodeListByID(c *gin.Context) {
	authContext := getAuthContext(c)
	nodeID, fcerr := extractNodeID(c)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get nodeID from request")
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	list, fcerr := r.managers.Node.ListByID(authContext, nodeID)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (r *Router) getNodeContentByID(c *gin.Context) {
	authContext := getAuthContext(c)
	nodeID, fcerr := extractNodeID(c)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to get nodeID from request")
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	node, reader, size, fcerr := r.managers.Node.DownloadFile(authContext, nodeID)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}
	defer reader.Close()

	c.DataFromReader(http.StatusOK, size, string(node.MimeType), reader, nil)
}

func (r *Router) createNodeByID(c *gin.Context) {
	authContext := getAuthContext(c)

	node := &models.Node{}
	err := c.BindJSON(node)
	if err != nil {
		fcerr := fcerror.NewError(fcerror.ErrBadRequest, err)
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	created, fcerr := r.managers.Node.CreateNode(authContext, node)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": node, "created": created})
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
		fcerr = fcerror.NewErrorSkipFunc(fcerror.ErrBadRequest, errors.New("NodeID not found in path param"))
		return
	}

	nodeID = models.NodeID(nodeIDStr)
	return
}
