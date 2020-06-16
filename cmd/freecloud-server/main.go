package main

import (
	"github.com/freecloudio/server/application"
	"github.com/freecloudio/server/application/persistence"

	_ "github.com/freecloudio/server/plugin/dgraph"
)

func main() {
	persistence.InitializeUsedPlugins()

	userMgr := application.UserManager{}
	userMgr.CreateUser()
}
