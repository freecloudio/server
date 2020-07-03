package fcerror

const (
	ErrIDTokenInvalid = iota + 200
	ErrIDNotAuthorized
)

func init() {
	errorDescriptions[ErrIDTokenInvalid] = "Token not valid"
	errorDescriptions[ErrIDNotAuthorized] = "Not authorized for this action"
}
