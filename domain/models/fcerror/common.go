package fcerror

const (
	ErrUnknown ErrorID = iota + 1
	ErrModelConversionFailed
	ErrBadRequest
	ErrInternalServerError
)

func init() {
	errorDescriptions[ErrUnknown] = "Unknown Error"
	errorDescriptions[ErrModelConversionFailed] = "Failed to convert models"
	errorDescriptions[ErrBadRequest] = "Bad Request"
	errorDescriptions[ErrInternalServerError] = "Internal Server Error"
}
