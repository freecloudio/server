package fcerror

const (
	ErrDBTransactionCreationFailed = iota + 300
	ErrDBWriteFailed
	ErrDBReadFailed
	ErrDBCommitFailed
	ErrDBRollbackFailed
	ErrDBAuthentication
	ErrDBUnavailable
	ErrDBCloseSessionFailed
	ErrDBInitializationFailed
	ErrDBCloseFailed
)

func init() {
	errorDescriptions[ErrDBTransactionCreationFailed] = "Failed to create a DB transaction"
	errorDescriptions[ErrDBWriteFailed] = "Failed to write data to DB"
	errorDescriptions[ErrDBReadFailed] = "Failed to read data from DB"
	errorDescriptions[ErrDBCommitFailed] = "Failed to commit DB transaction"
	errorDescriptions[ErrDBRollbackFailed] = "Failed to rollback DB transaction"
	errorDescriptions[ErrDBAuthentication] = "Failed to authenticate to the DB"
	errorDescriptions[ErrDBUnavailable] = "Database is unavailable"
	errorDescriptions[ErrDBCloseSessionFailed] = "Failed to close database session"
	errorDescriptions[ErrDBInitializationFailed] = "Failed to initialize database connection"
	errorDescriptions[ErrDBCloseFailed] = "Failed to close database connection"
}
