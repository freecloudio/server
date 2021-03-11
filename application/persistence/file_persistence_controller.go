package persistence

import (
	"github.com/freecloudio/server/domain/models/fcerror"
)

type FilePersistenceController interface {
	StartReadTransaction() (FilePersistenceReadTransaction, *fcerror.Error)
	StartReadWriteTransaction() (FilePersistenceReadWriteTransaction, *fcerror.Error)
}

type FilePersistenceReadTransaction interface {
	ReadTransaction
}

type FilePersistenceReadWriteTransaction interface {
	ReadWriteTransaction
	FilePersistenceReadTransaction
}
