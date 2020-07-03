package fcerror

const (
	ErrIDDBTransactionCreationFailed = iota + 300
	ErrIDDBWriteFailed
	ErrIDDBCommitFailed
	ErrIDDBRollbackFailed
	ErrIDDBAuthentication
	ErrIDDBUnavailable
)

func init() {
	errorDescriptions[ErrIDDBTransactionCreationFailed] = "Failed to create a DB transaction"
	errorDescriptions[ErrIDDBWriteFailed] = "Failed to write data to DB"
	errorDescriptions[ErrIDDBCommitFailed] = "Failed to commit DB transaction"
	errorDescriptions[ErrIDDBRollbackFailed] = "Failed to rollback DB transaction"
	errorDescriptions[ErrIDDBAuthentication] = "Failed to authenticate to the DB"
	errorDescriptions[ErrIDDBUnavailable] = "Database is unavailable"
}
