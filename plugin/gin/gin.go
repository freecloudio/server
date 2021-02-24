package gin

import (
	"net/http"

	"github.com/freecloudio/server/application"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type Router struct {
	engine *gin.Engine

	userMgr application.AuthManager
}

func NewRouter(userMgr application.AuthManager) (router *Router) {
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(ginlogrus.Logger(logrus.New()))

	router = &Router{
		engine:  ginRouter,
		userMgr: userMgr,
	}
	router.buildRoutes()

	return
}

func (r *Router) Serve(addr string) {
	r.engine.Run(addr)
}

func (r *Router) buildRoutes() {
	r.buildUserRoutes()

	r.engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	})
}

func errToStatus(fcerr *fcerror.Error) int {
	switch fcerr.ID {
	case fcerror.ErrUserNotFound:
		return http.StatusNotFound
	case fcerror.ErrBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
