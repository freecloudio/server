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
