package viperplg

import (
	"os"
	"time"

	"github.com/freecloudio/server/application/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	keyUserPersistencePlugin = "persistence.user.plugin"
	keyAuthPersistencePlugin = "persistence.auth.plugin"

	keyAuthTokenValueLength = "auth.token.length"
	keyAuthTokenExpiration  = "auth.token.expiration"
)

type ViperConfig struct {
	viper *viper.Viper
}

func InitViperConfig() *ViperConfig {
	v := viper.New()
	p := pflag.NewFlagSet("freecloud-server", pflag.ExitOnError)

	p.String(keyUserPersistencePlugin, string(config.NeoPersistenceKey), "Key of the persistence plugin to use for user management")
	p.String(keyAuthPersistencePlugin, string(config.NeoPersistenceKey), "Key of the persistence plugin to use for auth management")

	p.Int(keyAuthTokenValueLength, 32, "Length of the token used for authentication")
	p.Int(keyAuthTokenExpiration, 24, "Time a token is valid in hours")

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

func (cfg *ViperConfig) GetUserPersistencePluginKey() config.PersistencePluginKey {
	return config.PersistencePluginKey(cfg.viper.GetString(keyUserPersistencePlugin))
}

func (cfg *ViperConfig) GetAuthPersistencePluginKey() config.PersistencePluginKey {
	return config.PersistencePluginKey(cfg.viper.GetString(keyAuthPersistencePlugin))
}

func (cfg *ViperConfig) GetTokenValueLength() int {
	return cfg.viper.GetInt(keyAuthTokenValueLength)
}

func (cfg *ViperConfig) GetTokenExpirationDuration() time.Duration {
	return time.Duration(cfg.viper.GetInt(keyAuthTokenValueLength)) * time.Hour
}
