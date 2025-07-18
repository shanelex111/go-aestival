package token

import (
	"encoding/json"
	"errors"
	"sort"
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
		accountPrefixKey = getAccountPrefixKey(c.Account.ID)
		pipe             = redis.RDB.Pipeline()
	)

	// 查询所有
	existTokens, err := c.findAll()
	if err != nil {
		return err
	}
	if len(existTokens) > 0 {
		// 先删除device相同的token
		for _, v := range existTokens {
			if v.Device.DeviceID == c.Device.DeviceID &&
				v.Device.DeviceModel == c.Device.DeviceModel &&
				v.Device.DeviceType == c.Device.DeviceType &&
				v.Device.AppVersion == c.Device.AppVersion {
				pipe.Del(redis.Ctx, accessPrefix+v.Access.Token)
				pipe.Del(redis.Ctx, refreshPrefix+v.Access.Refresh)
				pipe.HDel(redis.Ctx, accountPrefixKey, v.Access.Refresh)
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

	// 限制token数量
	if err := c.limit(); err != nil {
		return err
	}
	return nil
}

func (c *CacheToken) limit() error {
	existTokens, err := c.findAll()
	if err != nil {
		return err
	}
	if len(existTokens) <= cfg.CacheConfig.Limit {
		return nil
	}

	var (
		pipe             = redis.RDB.Pipeline()
		accountPrefixKey = getAccountPrefixKey(c.Account.ID)
	)
	pipe.Del(redis.Ctx, accessPrefix+existTokens[0].Access.Token)
	pipe.Del(redis.Ctx, refreshPrefix+existTokens[0].Access.Refresh)
	pipe.HDel(redis.Ctx, accountPrefixKey, existTokens[0].Access.Refresh)
	_, err = pipe.Exec(redis.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheToken) findAll() ([]*CacheToken, error) {
	var (
		accountPrefixKey = getAccountPrefixKey(c.Account.ID)
	)
	result, err := redis.RDB.HGetAll(redis.Ctx, accountPrefixKey).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var tokens []*CacheToken
	for _, v := range result {
		var token CacheToken
		if err := json.Unmarshal([]byte(v), &token); err != nil {
			continue
		}
		tokens = append(tokens, &token)
	}

	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].Device.CreatedAt < tokens[j].Device.CreatedAt
	})
	return tokens, nil
}

func getAccountPrefixKey(accountID uint) string {
	return accountPrefix + strconv.FormatUint(uint64(accountID), 10)
}

func GetRefresh(refresh string) (*CacheToken, error) {
	val, err := redis.RDB.Get(redis.Ctx, refreshPrefix+refresh).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var token CacheToken
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (c *CacheToken) Delete() error {
	var (
		pipe             = redis.RDB.Pipeline()
		accountPrefixKey = getAccountPrefixKey(c.Account.ID)
	)

	pipe.Del(redis.Ctx, accessPrefix+c.Access.Token)
	pipe.Del(redis.Ctx, refreshPrefix+c.Access.Refresh)
	pipe.HDel(redis.Ctx, accountPrefixKey, c.Access.Refresh)
	_, err := pipe.Exec(redis.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func Find(token string) (*CacheToken, error) {
	val, err := redis.RDB.Get(redis.Ctx, accessPrefix+token).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var c CacheToken
	if err := json.Unmarshal([]byte(val), &c); err != nil {
		return nil, err
	}
	return &c, nil
}
func DelAllByAccountID(accountID uint) error {
	var (
		accountPrefixKey = getAccountPrefixKey(accountID)
	)

	// 获取所有token
	result, err := redis.RDB.HGetAll(redis.Ctx, accountPrefixKey).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil
		}
		return err
	}
	if len(result) > 0 {
		var pipe = redis.RDB.Pipeline()
		for _, v := range result {
			var token CacheToken
			if err := json.Unmarshal([]byte(v), &token); err != nil {
				continue
			}
			pipe.Del(redis.Ctx, accessPrefix+token.Access.Token)
			pipe.Del(redis.Ctx, refreshPrefix+token.Access.Refresh)
		}
		pipe.Del(redis.Ctx, accountPrefixKey)
		if _, err := pipe.Exec(redis.Ctx); err != nil {
			return err
		}
	}
	return nil
}
