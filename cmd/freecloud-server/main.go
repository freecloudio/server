package main

import (
	"github.com/freecloudio/server/application"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models"

	_ "github.com/freecloudio/server/plugin/neo"
)

func main() {
	persistence.InitializeUsedPlugins()

	userMgr := application.UserManager{}
	//userMgr.CreateUser()
	userMgr.GetUser(models.UserID(1))
}
