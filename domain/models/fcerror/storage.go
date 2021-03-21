package fcerror

const (
	ErrFileFolderCreationFailed ErrorID = iota + 500
	ErrStorageOperationWithWrongUserPerspective
	ErrOpenUploadFile
	ErrOpenUserFile
	ErrCopyFileFailed
)

func init() {
	errorDescriptions[ErrFileFolderCreationFailed] = "Failed to create new folder or file"
	errorDescriptions[ErrStorageOperationWithWrongUserPerspective] = "Server tried a storage operation from the wrong file user perspective"
	errorDescriptions[ErrOpenUploadFile] = "Failed to open uploaded file"
	errorDescriptions[ErrOpenUserFile] = "Failed to open users file"
	errorDescriptions[ErrCopyFileFailed] = "Failed to copy file"
}
