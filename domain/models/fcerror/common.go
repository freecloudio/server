package fcerror

const (
	ErrUnknown ErrorID = iota + 1
	ErrModelConversionFailed
	ErrBadRequest
)

func init() {
	errorDescriptions[ErrUnknown] = "Unknown Error"
	errorDescriptions[ErrModelConversionFailed] = "Failed to convert models"
	errorDescriptions[ErrBadRequest] = "Bad Request"
}
