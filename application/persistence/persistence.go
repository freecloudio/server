package persistence

import "github.com/sirupsen/logrus"

// PluginInitializationFunc is a func for initializing a persistence plugin
type PluginInitializationFunc func() error

// TODO: Proper type for persistence types
var pluginInitFuncs map[string]PluginInitializationFunc
var pluginsUsed map[string]struct{}

// RegisterPluginInitialization registers the init func of a persistence
// The func will only be called if the persitence is used
func RegisterPluginInitialization(persistencePluginKey string, initFunc PluginInitializationFunc) {
	if pluginInitFuncs == nil {
		pluginInitFuncs = make(map[string]PluginInitializationFunc)
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

func markPluginUsed(persistencePluginKey string) {
	if pluginsUsed == nil {
		pluginsUsed = make(map[string]struct{})
	}

	pluginsUsed[persistencePluginKey] = struct{}{}
}

// Transaction stands for a transaction of any persistence plugin
type Transaction interface {
	Commit() error
	Rollback() error
}
