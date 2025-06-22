package token

import (
	"encoding/json"
	"strconv"
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

	tokenBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	deviceBytes, err := json.Marshal(c.Device)
	if err != nil {
		return err
	}

	pipe := redis.RDB.Pipeline()
	pipe.Set(redis.Ctx,
		accessPrefix+c.Access.Token,
		tokenBytes,
		cfg.CacheConfig.AccessValid)

	pipe.Set(redis.Ctx,
		refreshPrefix+c.Access.Refresh,
		tokenBytes,
		cfg.CacheConfig.RefreshValid)

	pipe.HSet(redis.Ctx,
		accountPrefix+strconv.FormatUint(uint64(c.Account.ID), 10),
		c.Access.Token,
		deviceBytes,
	)
	_, err = pipe.Exec(redis.Ctx)
	if err != nil {
		return err
	}
	return nil
}
