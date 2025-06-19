package token

import (
	"encoding/json"

	"github.com/shanelex111/go-common/pkg/cache/redis"
)

func Create(c *CacheToken) error {
	valueBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	pipe := redis.RDB.Pipeline()
	pipe.Set(redis.Ctx,
		accessPrefix+c.Access.Token,
		valueBytes,
		cfg.CacheConfig.AccessExpired)

	pipe.Set(redis.Ctx,
		refresPrefix+c.Access.Refresh,
		valueBytes,
		cfg.CacheConfig.RefreshExpired)

	_, err = pipe.Exec(redis.Ctx)
	if err != nil {
		return err
	}
	return nil
}
