package viperplg

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	keyAuthSessionTokenLength     = "auth.session.token.length"
	keyAuthSessionExpiration      = "auth.session.expiration"
	keyAuthSessionCleanupInterval = "auth.session.cleanup.interval"

	keyDBConnectionUsername = "db.connection.username"
	keyDBConnectionPassword = "db.connection.password"
	keyDBConnectionString   = "db.connection.string"

	keyFileStorageTempBasePath    = "storage.temp.basepath"
	keyFileStorageLocalFSBasePath = "storage.file.localfs.basepath"
)

type ViperConfig struct {
	viper *viper.Viper
}

func InitViperConfig() *ViperConfig {
	v := viper.New()
	p := pflag.NewFlagSet("freecloud-server", pflag.ExitOnError)

	p.Int(keyAuthSessionTokenLength, 32, "Length of the token used for authentication")
	p.Int(keyAuthSessionExpiration, 24, "Time a session is valid in hours")
	p.Int(keyAuthSessionCleanupInterval, 1, "Interval in which expired sessions will be cleaned in hours")

	p.String(keyDBConnectionUsername, "neo4j", "Username for the database connection")
	p.String(keyDBConnectionPassword, "freecloud", "Password for the database connection")
	p.String(keyDBConnectionString, "bolt://localhost:7687", "Connection string for the database")

	p.String(keyFileStorageTempBasePath, "tmp", "Base path of folder for temporary files")
	p.String(keyFileStorageLocalFSBasePath, "data", "Base path of the local filesystem file storage")

	p.Parse(os.Args[1:])
	v.BindPFlags(p)
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		// Error is not a file not found error
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logrus.WithError(err).Fatal("Failed to read config file")
		}
	}

	return &ViperConfig{v}
}

func (cfg *ViperConfig) GetSessionTokenLength() int {
	return cfg.viper.GetInt(keyAuthSessionTokenLength)
}

func (cfg *ViperConfig) GetSessionExpirationDuration() time.Duration {
	return time.Duration(cfg.viper.GetInt(keyAuthSessionExpiration)) * time.Hour
}

func (cfg *ViperConfig) GetSessionCleanupInterval() time.Duration {
	return time.Duration(cfg.viper.GetInt(keyAuthSessionCleanupInterval)) * time.Hour
}

func (cfg *ViperConfig) GetDBUsername() string {
	return cfg.viper.GetString(keyDBConnectionUsername)
}

func (cfg *ViperConfig) GetDBPassword() string {
	return cfg.viper.GetString(keyDBConnectionPassword)
}

func (cfg *ViperConfig) GetDBConnectionString() string {
	return cfg.viper.GetString(keyDBConnectionString)
}

func (cfg *ViperConfig) GetFileStorageTempBasePath() string {
	return cfg.viper.GetString(keyFileStorageTempBasePath)
}

func (cfg *ViperConfig) GetFileStorageLocalFSBasePath() string {
	return cfg.viper.GetString(keyFileStorageLocalFSBasePath)
}
