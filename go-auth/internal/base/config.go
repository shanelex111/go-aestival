package base

type BaseConfigEntity struct {
	TableName string `mapstructure:"table_name"`
}

type BaseCacheConfig struct {
	Prefix string `mapstructure:"prefix"`
}
