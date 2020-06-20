package config

type PersistencePluginKey string

const (
	DGraphPersistenceKey = PersistencePluginKey("dgraph")
	NeoPersistenceKey    = PersistencePluginKey("Neo")
)

func GetUserPersistencePluginKey() PersistencePluginKey {
	return NeoPersistenceKey
}
