package main

import (
	"github.com/freecloudio/server/application"
	"github.com/freecloudio/server/application/persistence"

	"github.com/freecloudio/server/plugin/gin"
	_ "github.com/freecloudio/server/plugin/neo"
	"github.com/freecloudio/server/plugin/viperplg"
)

func main() {
	persistence.InitializeUsedPlugins()

	cfg := viperplg.InitViperConfig()

	authMgr := application.NewAuthManager(cfg)

	router := gin.NewRouter(authMgr)
	router.Serve(":8080")
}
