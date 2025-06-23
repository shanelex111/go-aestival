package token

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/shanelex111/go-common/pkg/cache/redis"
	"github.com/shanelex111/go-common/pkg/util"
)

func (c *CacheToken) Create() error {
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

	var (
		accountPrefixKey = c.getAccountPrefixKey()
		pipe             = redis.RDB.Pipeline()
	)

	// 查询所有
	existTokens, err := c.FindAll()
	if err != nil {
		return err
	}
	if len(existTokens) > 0 {
		// 先删除device相同的token
		for k, v := range existTokens {
			if v.Device.DeviceID == c.Device.DeviceID &&
				v.Device.DeviceModel == c.Device.DeviceModel &&
				v.Device.DeviceType == c.Device.DeviceType &&
				v.Device.AppVersion == c.Device.AppVersion {
				pipe.Del(redis.Ctx, accessPrefix+v.Access.Token)
				pipe.Del(redis.Ctx, refreshPrefix+k)
				pipe.HDel(redis.Ctx, accountPrefixKey, k)
			}
		}
	}
	// 新生成的token
	pipe.Set(redis.Ctx,
		accessPrefix+c.Access.Token,
		tokenBytes,
		cfg.CacheConfig.AccessValid)

	pipe.Set(redis.Ctx,
		refreshPrefix+c.Access.Refresh,
		tokenBytes,
		cfg.CacheConfig.RefreshValid)

	pipe.HSet(redis.Ctx,
		accountPrefixKey,
		c.Access.Refresh,
		tokenBytes,
	)
	pipe.HExpire(redis.Ctx, accountPrefixKey, cfg.CacheConfig.RefreshValid, c.Access.Refresh)
	pipe.Expire(redis.Ctx, accountPrefixKey, cfg.CacheConfig.RefreshValid)
	_, err = pipe.Exec(redis.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheToken) Limit() error {

	return nil
}

func (c *CacheToken) FindAll() (map[string]*CacheToken, error) {
	var (
		accountPrefixKey = c.getAccountPrefixKey()
	)
	result, err := redis.RDB.HGetAll(redis.Ctx, accountPrefixKey).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	tokens := make(map[string]*CacheToken)
	for k, v := range result {
		var token CacheToken
		if err := json.Unmarshal([]byte(v), &token); err != nil {
			continue
		}

		tokens[k] = &token
	}

	return tokens, nil
}

func (c *CacheToken) getAccountPrefixKey() string {
	return accountPrefix + strconv.FormatUint(uint64(c.Account.ID), 10)
}
