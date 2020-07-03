package fcerror

const (
	ErrIDUserNotFound = iota + 100
)

func init() {
	errorDescriptions[ErrIDUserNotFound] = "User not found"
}
