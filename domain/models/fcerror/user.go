package fcerror

const (
	ErrUserNotFound = iota + 100
	ErrEmailAlreadyRegistered
)

func init() {
	errorDescriptions[ErrUserNotFound] = "User not found"
	errorDescriptions[ErrEmailAlreadyRegistered] = "User with this email address is already registered"
}
