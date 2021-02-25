package fcerror

const (
	ErrTokenInvalid = iota + 200
	ErrUnauthorized
	ErrForbidden
	ErrPasswordHashingFailed
	ErrTokenNotFound
)

func init() {
	errorDescriptions[ErrTokenInvalid] = "Token not valid"
	errorDescriptions[ErrUnauthorized] = "Not authorized for this action"
	errorDescriptions[ErrForbidden] = "This action is forbidden"
	errorDescriptions[ErrPasswordHashingFailed] = "Failed to hash password or stored hash is invalid"
	errorDescriptions[ErrTokenNotFound] = "Token could not be found"
}
