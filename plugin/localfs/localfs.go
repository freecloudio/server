package localfs

import (
	"fmt"
	"io"
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
		fcerr = fcerror.NewError(fcerror.ErrFileFolderCreationFailed, err)
	}
	return
}

func (*LocalFSStorage) Close() *fcerror.Error {
	return nil
}

func (fs *LocalFSStorage) getUserFolder(userID models.UserID) string {
	return utils.JoinPaths(fs.basepath, fmt.Sprintf("%d", userID))
}

func (fs *LocalFSStorage) getUserNodePath(node *models.Node) string {
	return utils.JoinPaths(fs.getUserFolder(node.OwnerID), node.FullPath)
}

func (fs *LocalFSStorage) CreateUserRootFolder(userID models.UserID) (fcerr *fcerror.Error) {
	userPath := fs.getUserFolder(userID)
	err := os.Mkdir(userPath, osPermission)
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrFileFolderCreationFailed, err)
	}
	return
}

func (fs *LocalFSStorage) CreateEmptyFileOrFolder(node *models.Node) (fcerr *fcerror.Error) {
	if node.OwnerID != node.PerspectiveUserID {
		return fcerror.NewError(fcerror.ErrFileFolderCreationFailed, nil)
	}

	path := fs.getUserNodePath(node)

	var err error
	switch node.Type {
	case models.NodeTypeFolder:
		err = os.Mkdir(path, osPermission)
	default:
		var file *os.File
		file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, osPermission)
		_ = file.Close()
	}

	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrFileFolderCreationFailed, err)
	}
	return
}

func (fs *LocalFSStorage) CopyFileFromUpload(node *models.Node, uploadPath string) (fcerr *fcerror.Error) {
	if node.OwnerID != node.PerspectiveUserID {
		return fcerror.NewError(fcerror.ErrFileFolderCreationFailed, nil)
	}

	path := fs.getUserNodePath(node)

	source, err := os.OpenFile(uploadPath, os.O_RDONLY, osPermission)
	if err != nil {
		return fcerror.NewError(fcerror.ErrOpenUploadFile, err)
	}
	defer source.Close()

	destination, err := os.OpenFile(path, os.O_RDWR, osPermission)
	if err != nil {
		return fcerror.NewError(fcerror.ErrOpenUploadFile, err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fcerror.NewError(fcerror.ErrBadRequest, err)
	}

	return
}
