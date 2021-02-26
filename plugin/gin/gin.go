package gin

import (
	"context"
	"net/http"

	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type Router struct {
	engine  *gin.Engine
	authMgr manager.AuthManager
	srv     *http.Server
}

func NewRouter(authMgr manager.AuthManager, addr string) (router *Router) {
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(ginlogrus.Logger(logrus.New()))
	ginRouter.Use(getAuthMiddleware(authMgr))

	router = &Router{
		engine:  ginRouter,
		authMgr: authMgr,
		srv: &http.Server{
			Addr:    ":8080",
			Handler: ginRouter,
		},
	}
	router.buildRoutes()

	return
}

func (r *Router) Serve() error {
	return r.srv.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}

func (r *Router) buildRoutes() {
	r.buildAuthRoutes()
	r.buildUserRoutes()

	r.engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	})
}

func errToStatus(fcerr *fcerror.Error) int {
	switch fcerr.ID {
	case fcerror.ErrUnauthorized:
		return http.StatusUnauthorized
	case fcerror.ErrForbidden:
		return http.StatusForbidden
	case fcerror.ErrUserNotFound:
		return http.StatusNotFound
	case fcerror.ErrBadRequest, fcerror.ErrEmailAlreadyRegistered:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
