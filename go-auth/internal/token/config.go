package token

import (
	"time"

	"github.com/spf13/viper"
)

const (
	defaultKey = "token"
)

var (
	cfg          *config
	accessPrefix string
	refresPrefix string
)

type config struct {
	CacheConfig *cacheConfig `mapstructure:"cache"`
}

type cacheConfig struct {
	Prefix         string        `mapstructure:"prefix"`
	AccessExpired  time.Duration `mapstructure:"access_expired"`
	RefreshExpired time.Duration `mapstructure:"refresh_expired"`
}

func Load(v *viper.Viper) {
	initConfig(v)
}

func initConfig(v *viper.Viper) {
	cfg := &config{}
	if err := v.Sub(defaultKey).Unmarshal(cfg); err != nil {
		panic(err)
	}
	accessPrefix = cfg.CacheConfig.Prefix + "access:"
	refresPrefix = cfg.CacheConfig.Prefix + "refresh:"
}
