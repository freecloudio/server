package storage

import (
	"io"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
)

type FileStorageController interface {
	CreateUserRootFolder(userID models.UserID) *fcerror.Error
	CreateEmptyFileOrFolder(node *models.Node) *fcerror.Error
	CopyFileFromUpload(node *models.Node, uploadPath string) *fcerror.Error
	DownloadFile(node *models.Node) (io.ReadCloser, int64, *fcerror.Error)
}
