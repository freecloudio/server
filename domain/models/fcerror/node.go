package fcerror

const (
	ErrNodeNotFound ErrorID = iota + 400
)

func init() {
	errorDescriptions[ErrNodeNotFound] = "File or folder not found"
}
