package device

import "go-auth/internal/base"

type DeviceEntity struct {
	base.BaseModelEntity
}

func (DeviceEntity) TableName() string {
	return cfg.ConfigEntity.TableName
}
