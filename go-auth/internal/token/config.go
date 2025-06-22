package token

import (
	"time"

	"github.com/spf13/viper"
)

const (
	defaultKey = "token"
)

var (
	cfg           *config
	accessPrefix  string
	refreshPrefix string
	accountPrefix string
)

type config struct {
	CacheConfig *cacheConfig `mapstructure:"cache"`
}

type cacheConfig struct {
	Prefix       string        `mapstructure:"prefix"`
	AccessValid  time.Duration `mapstructure:"access_valid"`
	RefreshValid time.Duration `mapstructure:"refresh_valid"`
}

func Load(v *viper.Viper) {
	initConfig(v)
}

func initConfig(v *viper.Viper) {
	cfg = &config{}
	if err := v.Sub(defaultKey).Unmarshal(cfg); err != nil {
		panic(err)
	}
	accessPrefix = cfg.CacheConfig.Prefix + "access:"
	refreshPrefix = cfg.CacheConfig.Prefix + "refresh:"
	accountPrefix = cfg.CacheConfig.Prefix + "account:"
}
