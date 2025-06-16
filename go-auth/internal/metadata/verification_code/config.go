package verification_code

import (
	"go-auth/internal/base"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultKey = "verification_code"
)

var (
	cfg *config
)

type config struct {
	ConfigEntity *base.BaseConfigEntity `mapstructure:"entity"`
	Number       int                    `mapstructure:"number"`
	Limited      int                    `mapstructure:"limited"`
	Period       time.Duration          `mapstructure:"period"`
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

func GetLimited() int {
	return cfg.Limited
}

func GetNumber() int {
	return cfg.Number
}
