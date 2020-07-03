package fcerror

const (
	ErrIDUnknown = iota
	ErrIDModelConversionFailed
)

func init() {
	errorDescriptions[ErrIDUnknown] = "Unknown Error"
	errorDescriptions[ErrIDModelConversionFailed] = "Failed to convert models"
}
