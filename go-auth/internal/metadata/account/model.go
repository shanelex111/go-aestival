package account

import "go-auth/internal/base"

const (
	StatusDeleted = -1
	StatusDisable = 0
	StatusEnable  = 1
)

type AccountEntity struct {
	base.BaseModelEntity
	Status   int    `gorm:"column:status;type:tinyint not null;default:0;comment:状态：1：enable，0：disable，-1：deleted"`
	Email    string `gorm:"column:email;type:varchar(255) not null;comment:邮箱"`
	Password string `gorm:"column:password;type:varchar(255) not null;comment:密码"`
}

func (AccountEntity) TableName() string {
	return cfg.ConfigEntity.TableName
}
