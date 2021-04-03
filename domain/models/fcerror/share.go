package fcerror

const (
	ErrShareContainsOtherShares ErrorID = iota + 600
	ErrShareIsInsideShare
)

func init() {
	errorDescriptions[ErrShareContainsOtherShares] = "Folder that should be shared, contains other shared files"
	errorDescriptions[ErrShareIsInsideShare] = "Folder that should be shared, is inside a folder shared with you"
}
