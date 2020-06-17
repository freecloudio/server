package persistence

import (
	"github.com/freecloudio/server/config"

	"github.com/sirupsen/logrus"
)

// PluginInitializationFunc is a func for initializing a persistence plugin
type PluginInitializationFunc func() error

var pluginInitFuncs map[config.PersistencePluginKey]PluginInitializationFunc
var pluginsUsed map[config.PersistencePluginKey]struct{}

// RegisterPluginInitialization registers the init func of a persistence
// The func will only be called if the persitence is used
func RegisterPluginInitialization(persistencePluginKey config.PersistencePluginKey, initFunc PluginInitializationFunc) {
	if pluginInitFuncs == nil {
		pluginInitFuncs = make(map[config.PersistencePluginKey]PluginInitializationFunc)
	}

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
	if pluginsUsed == nil {
		pluginsUsed = make(map[config.PersistencePluginKey]struct{})
	}

	pluginsUsed[persistencePluginKey] = struct{}{}
}

// ReadWriteTransaction stands for a transaction of any persistence plugin
type ReadWriteTransaction interface {
	Commit() error
	Rollback() error
}
