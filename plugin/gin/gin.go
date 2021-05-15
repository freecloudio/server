package gin

import (
	"context"
	"net/http"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"

	"github.com/gin-gonic/gin"
	ginlogrus "github.com/toorop/gin-logrus"
)

type Router struct {
	engine   *gin.Engine
	managers *manager.Managers
	srv      *http.Server
	cfg      config.Config
	logger   utils.Logger
}

func NewRouter(managers *manager.Managers, cfg config.Config, addr string) (router *Router) {
	logger := utils.CreateLogger(cfg.GetLoggingConfig())

	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(ginlogrus.Logger(logger))
	ginRouter.Use(getAuthMiddleware(managers.Auth))

	router = &Router{
		engine:   ginRouter,
		managers: managers,
		srv: &http.Server{
			Addr:    ":8080",
			Handler: ginRouter,
		},
		cfg:    cfg,
		logger: logger,
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
	r.buildNodeRoutes()
	r.buildGraphQLRoutes()

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
