package gin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/domain/models"
)

const userIDParam = "user_id"

func (r *Router) buildUserRoutes() {
	r.engine.GET("/api/user/:"+userIDParam, r.getUserID)
}

func (r *Router) getUserID(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		logrus.WithError(err).Error("Failed to get userID from request")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, fcerr := r.userMgr.GetUser(authorization.NewSystem(), userID)
	if err != nil {
		logrus.WithError(fcerr).Error("Failed to get user")
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusFound, user)
}

func extractUserID(c *gin.Context) (userID models.UserID, err error) {
	userIDStr := c.Param(userIDParam)
	if userIDStr == "" {
		err = fmt.Errorf("UserID not found in path params")
		return
	}

	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		return
	}

	userID = models.UserID(userIDInt)
	return
}
