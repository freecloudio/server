package gin

import (
	"context"
	"net/http"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/domain/models/fcerror"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type Router struct {
	engine   *gin.Engine
	managers *manager.Managers
	srv      *http.Server
	cfg      config.Config
}

func NewRouter(managers *manager.Managers, cfg config.Config, addr string) (router *Router) {
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(ginlogrus.Logger(logrus.New()))
	ginRouter.Use(getAuthMiddleware(managers.Auth))

	router = &Router{
		engine:   ginRouter,
		managers: managers,
		srv: &http.Server{
			Addr:    ":8080",
			Handler: ginRouter,
		},
		cfg: cfg,
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
	r.buildNodeRoutes()
	r.buildShareRoutes()

	r.engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	})
}

func errToStatus(fcerr *fcerror.Error) int {
	switch fcerr.ID {
	case fcerror.ErrUnauthorized, fcerror.ErrTokenNotFound:
		return http.StatusUnauthorized
	case fcerror.ErrForbidden:
		return http.StatusForbidden
	case fcerror.ErrUserNotFound, fcerror.ErrNodeNotFound:
		return http.StatusNotFound
	case fcerror.ErrBadRequest, fcerror.ErrEmailAlreadyRegistered:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
