package persistence

import (
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/domain/models/fcerror"

	"github.com/sirupsen/logrus"
)

// TODO: Plugin Close

// PluginInitializationFunc is a func for initializing a persistence plugin
type PluginInitializationFunc func() error // TODO: fcerror

var pluginInitFuncs = map[config.PersistencePluginKey]PluginInitializationFunc{}
var pluginsUsed = map[config.PersistencePluginKey]struct{}{}

// RegisterPluginInitialization registers the init func of a persistence
// The func will only be called if the persitence is used
func RegisterPluginInitialization(persistencePluginKey config.PersistencePluginKey, initFunc PluginInitializationFunc) {
	pluginInitFuncs[persistencePluginKey] = initFunc
}

// InitializeUsedPlugins call init funcs of all used persistence plugins
// If an error is returned, the server should be stopped
func InitializeUsedPlugins() (err error) {
	for key := range pluginsUsed {
		if initFunc, ok := pluginInitFuncs[key]; ok {
			err = initFunc()
			if err != nil {
				logrus.WithError(err).WithField("PersistencePluginKey", key).Error("Failed to init persistence plugin")
				return err
			}
		} else {
			logrus.WithField("PersistencePluginKey", key).Warn("No init func for used persistence plugin")
		}
	}

	return
}

func markPluginUsed(persistencePluginKey config.PersistencePluginKey) {
	pluginsUsed[persistencePluginKey] = struct{}{}
}

type ReadTransaction interface {
	Close() *fcerror.Error
}

// ReadWriteTransaction stands for a transaction of any persistence plugin
type ReadWriteTransaction interface {
	ReadTransaction
	Commit() *fcerror.Error
	Rollback() *fcerror.Error
}
