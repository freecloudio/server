package persistence

import (
	"github.com/freecloudio/server/domain/models/fcerror"
)

type ReadTransaction interface {
	Close() *fcerror.Error
}

// ReadWriteTransaction stands for a transaction of any persistence plugin
type ReadWriteTransaction interface {
	// Should either rollback or commit depending on preceding error
	Finish(fcerr *fcerror.Error) *fcerror.Error
	Commit() *fcerror.Error
	Rollback()
}
