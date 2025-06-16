package verification_code

import (
	"errors"
	"go-auth/internal/base"

	"github.com/shanelex111/go-common/pkg/cache/redis"
	"github.com/shanelex111/go-common/pkg/db/mysql"
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

func (e *VerificationCodeEntity) CountInEntity() (int64, error) {
	var (
		condition = &VerificationCodeEntity{
			Type:   e.Type,
			Target: e.Target,
		}
		count int64
	)

	if e.Type == base.SendCodeTypePhone {
		condition.CountryCode = e.CountryCode
	}

	if err := mysql.DB.Model(&VerificationCodeEntity{}).Where(condition).Where("deleted_at = 0").Count(&count).Error; err != nil {
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
