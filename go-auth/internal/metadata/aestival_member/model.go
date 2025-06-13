package aestival_member

import "go-auth/internal/base"

type AestivalMemberEntity struct {
	base.BaseModelEntity
}

func (AestivalMemberEntity) TableName() string {
	return cfg.ConfigEntity.TableName
}
