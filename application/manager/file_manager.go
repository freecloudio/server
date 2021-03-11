package manager

import (
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
)

type FileManager interface {
	Close()
}

func NewFileManager(cfg config.Config, filePersistence persistence.FilePersistenceController, managers *Managers) FileManager {
	fileMgr := &fileManager{
		cfg:             cfg,
		filePersistence: filePersistence,
		managers:        managers,
	}

	managers.File = fileMgr
	return fileMgr
}

type fileManager struct {
	cfg             config.Config
	filePersistence persistence.FilePersistenceController
	managers        *Managers
}

func (mgr *fileManager) Close() {
}
