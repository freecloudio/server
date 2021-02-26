package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freecloudio/server/application/manager"
	"github.com/freecloudio/server/application/persistence"
	"github.com/sirupsen/logrus"

	"github.com/freecloudio/server/plugin/gin"
	_ "github.com/freecloudio/server/plugin/neo"
	"github.com/freecloudio/server/plugin/viperplg"
)

func main() {
	fcerr := persistence.InitializeUsedPlugins()
	if fcerr != nil {
		logrus.WithError(fcerr).Fatal("Failed to initialize a persistence plugin - abort")
	}

	cfg := viperplg.InitViperConfig()

	authMgr := manager.NewAuthManager(cfg)

	router := gin.NewRouter(authMgr, ":8080")

	go func() {
		if err := router.Serve(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Failed start server - abort")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := router.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("Server forced to shutdown")
	}

	persistence.CloseUsedPlugins()
}
