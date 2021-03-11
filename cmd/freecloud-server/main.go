package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freecloudio/server/application/manager"
	"github.com/sirupsen/logrus"

	"github.com/freecloudio/server/plugin/gin"
	"github.com/freecloudio/server/plugin/neo"
	_ "github.com/freecloudio/server/plugin/neo"
	"github.com/freecloudio/server/plugin/viperplg"
)

func main() {
	authPersistence, fcerr := neo.CreateAuthPersistence()
	if fcerr != nil {
		logrus.WithError(fcerr).Fatal("Failed to initialize neo auth persistence plugin - abort")
	}
	userPersistence, fcerr := neo.CreateUserPersistence()
	if fcerr != nil {
		logrus.WithError(fcerr).Fatal("Failed to initialize neo auth persistence plugin - abort")
	}

	cfg := viperplg.InitViperConfig()

	authMgr := manager.NewAuthManager(cfg, userPersistence, authPersistence)
	userMgr := manager.NewUserManager(cfg, userPersistence)

	router := gin.NewRouter(authMgr, userMgr, ":8080")

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

	userMgr.Close()
	authMgr.Close()

	fcerr = userPersistence.Close()
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to close neo auth persistence plugin")
	}
	fcerr = authPersistence.Close()
	if fcerr != nil {
		logrus.WithError(fcerr).Fatal("Failed to close neo auth persistence plugin")
	}
}
