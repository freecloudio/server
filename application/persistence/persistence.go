package persistence

import (
	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/domain/models/fcerror"

	"github.com/sirupsen/logrus"
)

// PluginLifecycleFunc is a func for initializing or closing a persistence plugin
type PluginLifecycleFunc func() *fcerror.Error

type PluginLifecycleFuncs struct {
	InitializationFunc PluginLifecycleFunc
	CloseFunc          PluginLifecycleFunc
}

var pluginLifecycleFuncMap = map[config.PersistencePluginKey]PluginLifecycleFuncs{}

// Bool indicates whether plugin has been initialized
var pluginsUsed = map[config.PersistencePluginKey]bool{}

// RegisterPluginInitialization registers the init func of a persistence
// The func will only be called if the persitence is used
func RegisterPluginInitialization(persistencePluginKey config.PersistencePluginKey, lifecycleFuncs PluginLifecycleFuncs) {
	pluginLifecycleFuncMap[persistencePluginKey] = lifecycleFuncs
}

// InitializeUsedPlugins call init funcs of all used persistence plugins
// If an error is returned, the server should be stopped
func InitializeUsedPlugins() (fcerr *fcerror.Error) {
	for key := range pluginsUsed {
		if lifecycleFuncs, ok := pluginLifecycleFuncMap[key]; ok && lifecycleFuncs.InitializationFunc != nil {
			fcerr = lifecycleFuncs.InitializationFunc()
			if fcerr != nil {
				logrus.WithError(fcerr).WithField("PersistencePluginKey", key).Error("Failed to init persistence plugin")
				return
			}
			pluginsUsed[key] = true
		}
	}

	return
}

func CloseUsedPlugins() {
	for key, initialized := range pluginsUsed {
		if !initialized {
			continue
		}

		if lifecycleFuncs, ok := pluginLifecycleFuncMap[key]; ok && lifecycleFuncs.CloseFunc != nil {
			fcerr := lifecycleFuncs.CloseFunc()
			if fcerr != nil {
				logrus.WithError(fcerr).WithField("PersistencePluginKey", key).Error("Failed to close persistence plugin - Ignore for now")
			}
		}
	}
}

func markPluginUsed(persistencePluginKey config.PersistencePluginKey) {
	pluginsUsed[persistencePluginKey] = false
}

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
