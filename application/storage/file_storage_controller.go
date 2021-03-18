package storage

import (
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
)

type FileStorageController interface {
	CreateUserRootFolder(userID models.UserID) *fcerror.Error
	CreateEmptyFile(userID models.UserID, node *models.Node) *fcerror.Error
	CreateFileFromUpload(userID models.UserID, node *models.Node, uploadPath string) *fcerror.Error
}
