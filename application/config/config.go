package config

import "time"

type PersistencePluginKey string

const (
	NeoPersistenceKey = PersistencePluginKey("Neo")
)

type Config interface {
	GetUserPersistencePluginKey() PersistencePluginKey
	GetAuthPersistencePluginKey() PersistencePluginKey
	GetSessionTokenLength() int
	GetSessionExpirationDuration() time.Duration
	GetSessionCleanupInterval() time.Duration
}
