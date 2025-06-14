package account

import "go-auth/internal/base"

const (
	statusDeleted = -1
	statusDisable = 0
	statusEnable  = 1
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
