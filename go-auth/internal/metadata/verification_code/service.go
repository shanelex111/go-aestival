package verification_code

import (
	"errors"
	"go-auth/internal/base"

	goredis "github.com/redis/go-redis/v9"
	"github.com/shanelex111/go-common/pkg/cache/redis"
	"github.com/shanelex111/go-common/pkg/db/mysql"
	"github.com/shanelex111/go-common/pkg/util"
	"gorm.io/gorm"
)

func (e *VerificationCodeEntity) FindInEntity() (*VerificationCodeEntity, error) {
	var (
		entity    VerificationCodeEntity
		condition = &VerificationCodeEntity{
			Type:   e.Type,
			Status: e.Status,
			Target: e.Target,
			Code:   e.Code,
		}
	)
	if e.Type == base.SendCodeTypePhone {
		condition.CountryCode = e.CountryCode
	}

	if err := mysql.DB.Where(condition).Where("deleted_at = 0").Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (e *VerificationCodeEntity) FindLastInEntity() (*VerificationCodeEntity, error) {
	var (
		entity    VerificationCodeEntity
		condition = &VerificationCodeEntity{
			Type:   e.Type,
			Target: e.Target,
			Status: e.Status,
		}
	)
	if e.Type == base.SendCodeTypePhone {
		condition.CountryCode = e.CountryCode
	}

	if err := mysql.DB.Where(condition).Where("deleted_at = 0").Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (e *VerificationCodeEntity) CountTodayInEntity() (int64, error) {
	var (
		condition = &VerificationCodeEntity{
			Type:   e.Type,
			Target: e.Target,
		}
		count                int64
		todayStart, todayEnd = util.GetTodayMilli()
	)

	if e.Type == base.SendCodeTypePhone {
		condition.CountryCode = e.CountryCode
	}

	if err := mysql.DB.Model(&VerificationCodeEntity{}).
		Where(condition).
		Where("deleted_at = 0").
		Where("created_at >= ? AND created_at <= ?", todayStart, todayEnd).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (e *VerificationCodeEntity) ExpiredAllInEntity() error {
	var (
		condition = &VerificationCodeEntity{
			Type:   e.Type,
			Target: e.Target,
			Status: e.Status,
		}
	)
	if e.Type == base.SendCodeTypePhone {
		condition.CountryCode = e.CountryCode
	}

	return mysql.DB.Model(&VerificationCodeEntity{}).Where(condition).Where("deleted_at = 0").Update("status", StatusExpired).Error
}

func (e *VerificationCodeEntity) SaveInEntity() error {
	return mysql.DB.Save(e).Error
}

func (e *VerificationCodeEntity) SaveInCache() error {
	var (
		redisKey = cfg.Cache.Prefix + e.Type + ":"
	)

	if e.Type == base.SendCodeTypePhone {
		redisKey += e.CountryCode + ":" + e.Target
	} else {
		redisKey += e.Target
	}
	return redis.RDB.Set(redis.Ctx, redisKey, e.Code, cfg.Period).Err()
}

func (e *VerificationCodeEntity) DeleteInCache() error {
	var (
		redisKey = cfg.Cache.Prefix + e.Type + ":"
	)

	if e.Type == base.SendCodeTypePhone {
		redisKey += e.CountryCode + ":" + e.Target
	} else {
		redisKey += e.Target
	}
	return redis.RDB.Del(redis.Ctx, redisKey).Err()
}

func (e *VerificationCodeEntity) FindInCache() (string, error) {
	var (
		redisKey = cfg.Cache.Prefix + e.Type + ":"
	)

	if e.Type == base.SendCodeTypePhone {
		redisKey += e.CountryCode + ":" + e.Target
	} else {
		redisKey += e.Target
	}
	result, err := redis.RDB.Get(redis.Ctx, redisKey).Result()
	if err != nil {
		if !errors.Is(err, goredis.Nil) {
			return "", err
		}
	}
	return result, nil
}
