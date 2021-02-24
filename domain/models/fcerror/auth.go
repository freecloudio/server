package fcerror

const (
	ErrTokenInvalid = iota + 200
	ErrUnauthorized
	ErrForbidden
	ErrPasswordHashInvalid
	ErrPasswordHasingFailed
)

func init() {
	errorDescriptions[ErrTokenInvalid] = "Token not valid"
	errorDescriptions[ErrUnauthorized] = "Not authorized for this action"
	errorDescriptions[ErrForbidden] = "This action is forbidden"
	errorDescriptions[ErrPasswordHashInvalid] = "Password hash in database not valid"
	errorDescriptions[ErrPasswordHasingFailed] = "Failed to hash password"
}
