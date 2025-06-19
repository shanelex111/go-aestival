package token

import (
	"encoding/json"
	"time"

	"github.com/shanelex111/go-common/pkg/cache/redis"
	"github.com/shanelex111/go-common/pkg/util"
)

func Create(c *CacheToken) error {
	var (
		now = time.Now().UnixMilli()
	)

	c.Access = &CacheTokenAccess{
		Token:     util.GetUUID(),
		ExpiredAt: now + cfg.CacheConfig.AccessValid.Milliseconds(),

		Refresh:          util.GetUUID(),
		RefreshExpiredAt: now + cfg.CacheConfig.RefreshValid.Milliseconds(),
	}

	valueBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	pipe := redis.RDB.Pipeline()
	pipe.Set(redis.Ctx,
		accessPrefix+c.Access.Token,
		valueBytes,
		cfg.CacheConfig.AccessValid)

	pipe.Set(redis.Ctx,
		refreshPrefix+c.Access.Refresh,
		valueBytes,
		cfg.CacheConfig.RefreshValid)

	_, err = pipe.Exec(redis.Ctx)
	if err != nil {
		return err
	}
	return nil
}
