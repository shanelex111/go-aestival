package account

import "go-auth/internal/base"

const (
	StatusDeleted = -1
	StatusDisable = 0
	StatusEnable  = 1
)

type AccountEntity struct {
	base.BaseModelEntity
	Email            string `gorm:"column:email;type:varchar(255) not null;comment:邮箱"`
	PhoneCountryCode string `gorm:"column:phone_country_code;type:varchar(8) not null;comment:手机国家代码"`
	PhoneNumber      string `gorm:"column:phone_number;type:varchar(20) not null;comment:手机号码"`
	Password         string `gorm:"column:password;type:varchar(255) not null;comment:密码"`
	Status           int    `gorm:"column:status;type:tinyint not null;default:0;comment:状态：1：enable，0：disable，-1：deleted"`
}

func (AccountEntity) TableName() string {
	return cfg.ConfigEntity.TableName
}
