package config

type PersistencePluginKey string

const (
	NeoPersistenceKey = PersistencePluginKey("Neo")
)

func GetUserPersistencePluginKey() PersistencePluginKey {
	return NeoPersistenceKey
}

func GetAuthPersistencePluginKey() PersistencePluginKey {
	return NeoPersistenceKey
}
