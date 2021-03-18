package localfs

import (
	"fmt"
	"os"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"
)

type LocalFSStorage struct {
	basepath string
}

const osPermission os.FileMode = 0770

func CreateLocalFSStorage(cfg config.Config) (localFS *LocalFSStorage, fcerr *fcerror.Error) {
	localFS = &LocalFSStorage{
		basepath: cfg.GetFileStorageLocalFSBasePath(),
	}
	err := os.MkdirAll(localFS.basepath, osPermission)
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrFolderCreationFailed, err)
	}
	return
}

func (*LocalFSStorage) Close() *fcerror.Error {
	return nil
}

func (fs *LocalFSStorage) getUserFolder(userID models.UserID) (path string) {
	return utils.JoinPaths(fs.basepath, fmt.Sprintf("%d", userID))
}

func (fs *LocalFSStorage) CreateUserRootFolder(userID models.UserID) (fcerr *fcerror.Error) {
	userPath := fs.getUserFolder(userID)
	err := os.Mkdir(userPath, osPermission)
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrFolderCreationFailed, err)
	}
	return
}

func (fs *LocalFSStorage) CreateEmptyFile(userID models.UserID, node *models.Node) (fcerr *fcerror.Error) {
	return
}

func (fs *LocalFSStorage) CreateFileFromUpload(userID models.UserID, node *models.Node, uploadPath string) (fcerr *fcerror.Error) {
	return
}
