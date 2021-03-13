package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) buildNodeRoutes() {
	grp := r.engine.Group("/api/node")

	grp.GET("info/path/*path", r.getNodeInfoByPath)
	// info/id/
	// list/id/
	// content/id/
}

func (r *Router) getNodeInfoByPath(c *gin.Context) {
	authContext := getAuthContext(c)
	path := c.Param("path")

	node, fcerr := r.managers.Node.GetNodeByPath(authContext, authContext.User.ID, path)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, node)
}
