package fcerror

const (
	ErrFolderCreationFailed ErrorID = iota + 500
)

func init() {
	errorDescriptions[ErrFolderCreationFailed] = "Failed to create new folder"
}
