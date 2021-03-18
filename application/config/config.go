package config

import "time"

type PersistencePluginKey string

const (
	NeoPersistenceKey = PersistencePluginKey("Neo")
)

type Config interface {
	GetSessionTokenLength() int
	GetSessionExpirationDuration() time.Duration
	GetSessionCleanupInterval() time.Duration

	GetDBUsername() string
	GetDBPassword() string
	GetDBConnectionString() string

	GetFileStorageTempBasePath() string
	GetFileStorageLocalFSBasePath() string
}
