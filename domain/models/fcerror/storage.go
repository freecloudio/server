package fcerror

const (
	ErrFileFolderCreationFailed ErrorID = iota + 500
	ErrStorageOperationWithWrongUserPerspective
)

func init() {
	errorDescriptions[ErrFileFolderCreationFailed] = "Failed to create new folder or file"
	errorDescriptions[ErrStorageOperationWithWrongUserPerspective] = "Server tried a storage operation from the wrong file user perspective"
}
