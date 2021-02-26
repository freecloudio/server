package gin

import (
	"net/http"

	"github.com/freecloudio/server/application/authorization"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (r *Router) buildAuthRoutes() {
	grp := r.engine.Group("/api/auth")

	grp.POST("login", r.login)
	grp.POST("logout", r.logout)
}

func (r *Router) login(c *gin.Context) {
	user := &models.User{}
	err := c.BindJSON(user)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse request into user")
		fcerr := fcerror.NewError(fcerror.ErrBadRequest, err)
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	token, fcerr := r.authMgr.Login(user.Email, user.Password)
	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (r *Router) logout(c *gin.Context) {
	authContext := getAuthContext(c)

	var fcerr *fcerror.Error
	if tokenInt, ok := c.Get(authTokenKey); authContext.Type == authorization.ContextTypeUser && ok {
		token := tokenInt.(models.Token)
		fcerr = r.authMgr.Logout(token)
	} else {
		fcerr = fcerror.NewError(fcerror.ErrUnauthorized, nil)
	}

	if fcerr != nil {
		c.JSON(errToStatus(fcerr), fcerr)
		return
	}

	c.Status(http.StatusNoContent)
}
