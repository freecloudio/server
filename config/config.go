package config

type PersistencePluginKey string

const (
	DGraphPersistenceKey = PersistencePluginKey("dgraph")
)

func GetUserPersistencePluginKey() PersistencePluginKey {
	return DGraphPersistenceKey
}
