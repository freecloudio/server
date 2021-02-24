package fcerror

const (
	ErrDBTransactionCreationFailed = iota + 300
	ErrDBWriteFailed
	ErrDBCommitFailed
	ErrDBRollbackFailed
	ErrDBAuthentication
	ErrDBUnavailable
)

func init() {
	errorDescriptions[ErrDBTransactionCreationFailed] = "Failed to create a DB transaction"
	errorDescriptions[ErrDBWriteFailed] = "Failed to write data to DB"
	errorDescriptions[ErrDBCommitFailed] = "Failed to commit DB transaction"
	errorDescriptions[ErrDBRollbackFailed] = "Failed to rollback DB transaction"
	errorDescriptions[ErrDBAuthentication] = "Failed to authenticate to the DB"
	errorDescriptions[ErrDBUnavailable] = "Database is unavailable"
}
