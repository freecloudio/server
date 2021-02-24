package fcerror

const (
	ErrUserNotFound = iota + 100
	ErrEmailAlreadyRegisters
)

func init() {
	errorDescriptions[ErrUserNotFound] = "User not found"
	errorDescriptions[ErrEmailAlreadyRegisters] = "User with this email address is already registered"
}
