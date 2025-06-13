package aestival_member

import (
	"go-auth/internal/base"

	"github.com/spf13/viper"
)

const (
	defaultKey = "aestival_member"
)

var (
	cfg *config
)

type config struct {
	ConfigEntity *base.BaseConfigEntity `mapstructure:"entity"`
}

func Load(v *viper.Viper) {
	initConfig(v)
}

func initConfig(v *viper.Viper) {
	cfg = &config{
		ConfigEntity: &base.BaseConfigEntity{
			TableName: defaultKey,
		},
	}

	if v.IsSet(defaultKey) {
		if err := v.Sub(defaultKey).Unmarshal(cfg); err != nil {
			panic(err)
		}
	}

}
