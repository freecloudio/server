package fcerror

const (
	ErrShareContainsOtherShares ErrorID = iota + 600
)

func init() {
	errorDescriptions[ErrShareContainsOtherShares] = "Node that should be shared, contains other shared files"
}
