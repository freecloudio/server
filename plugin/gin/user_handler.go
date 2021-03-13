package gin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const userIDParam = "user_id"

func (r *Router) buildUserRoutes() {
	grp := r.engine.Group("/api/user")

	grp.GET("", r.getOwnUser)
	grp.POST("", r.registerUser)
	grp.GET(":"+userIDParam, r.getUserByID)
}

func (r *Router) registerUser(c *gin.Context) {
	authContext := getAuthContext(c)

	user := &models.User{}
	err := c.BindJSON(user)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse request into user")
		fcerr := fcerror.NewError(fcerror.ErrBadRequest, err)
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	session, fcerr := r.managers.User.CreateUser(authContext, user)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (r *Router) getOwnUser(c *gin.Context) {
	authContext := getAuthContext(c)

	if authContext.User != nil {
		c.JSON(http.StatusOK, authContext.User)
		return
	}

	fcerr := fcerror.NewError(fcerror.ErrUnauthorized, nil)
	c.JSON(errToStatus(fcerr), fcerr)
}

func (r *Router) getUserByID(c *gin.Context) {
	authContext := getAuthContext(c)

	userID, err := extractUserID(c)
	if err != nil {
		logrus.WithError(err).Error("Failed to get userID from request")
		fcerr := fcerror.NewError(fcerror.ErrBadRequest, err)
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	user, fcerr := r.managers.User.GetUserByID(authContext, userID)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, user)
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
