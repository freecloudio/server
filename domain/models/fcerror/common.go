package fcerror

const (
	ErrUnknown = iota
	ErrModelConversionFailed
	ErrBadRequest
)

func init() {
	errorDescriptions[ErrUnknown] = "Unknown Error"
	errorDescriptions[ErrModelConversionFailed] = "Failed to convert models"
	errorDescriptions[ErrBadRequest] = "Bad Request"
}
