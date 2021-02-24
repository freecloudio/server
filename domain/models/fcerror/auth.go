package fcerror

const (
	ErrTokenInvalid = iota + 200
	ErrNotAuthorized
	ErrPasswordHashInvalid
	ErrPasswordHasingFailed
)

func init() {
	errorDescriptions[ErrTokenInvalid] = "Token not valid"
	errorDescriptions[ErrNotAuthorized] = "Not authorized for this action"
	errorDescriptions[ErrPasswordHashInvalid] = "Password hash in database not valid"
	errorDescriptions[ErrPasswordHasingFailed] = "Failed to hash password"
}
