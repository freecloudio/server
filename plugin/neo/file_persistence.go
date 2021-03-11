package neo

import (
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/domain/models/fcerror"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

func init() {
}

type FilePersistence struct{}

func CreateFilePersistence(cfg config.Config) (filePersistence *FilePersistence, fcerr *fcerror.Error) {
	if neo == nil {
		fcerr = initializeNeo(cfg)
		if fcerr != nil {
			return
		}
	}
	filePersistence = &FilePersistence{}
	return
}

func (*FilePersistence) Close() *fcerror.Error {
	if neo != nil {
		return closeNeo()
	}
	return nil
}

func (*FilePersistence) StartReadTransaction() (tx persistence.FilePersistenceReadTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeRead)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo read transaction")
		return
	}
	return &fileReadTransaction{txCtx}, nil
}

func (*FilePersistence) StartReadWriteTransaction() (tx persistence.FilePersistenceReadWriteTransaction, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeWrite)
	if fcerr != nil {
		logrus.WithError(fcerr).Error("Failed to create neo write transaction")
		return
	}
	return &fileReadWriteTransaction{fileReadTransaction{txCtx}}, nil
}

type fileReadTransaction struct {
	*transactionCtx
}

type fileReadWriteTransaction struct {
	fileReadTransaction
}
