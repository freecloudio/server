package main

import (
	"github.com/freecloudio/server/application"
	"github.com/freecloudio/server/application/persistence"

	"github.com/freecloudio/server/plugin/gin"
	_ "github.com/freecloudio/server/plugin/neo"
)

func main() {
	persistence.InitializeUsedPlugins()

	userMgr := application.UserManager{}

	router := gin.NewRouter(&userMgr)
	router.Serve(":8080")
}
