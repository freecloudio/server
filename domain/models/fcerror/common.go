package fcerror

const (
	ErrIDUnknown = iota
	ErrIDModelConversionFailed
	ErrBadRequest
)

func init() {
	errorDescriptions[ErrIDUnknown] = "Unknown Error"
	errorDescriptions[ErrIDModelConversionFailed] = "Failed to convert models"
	errorDescriptions[ErrBadRequest] = "Bad Request"
}
