package account

import "go-auth/internal/base"

const (
	StatusDeleted = -1
	StatusDisable = 0
	StatusEnable  = 1
)

type Entity struct {
	base.BaseModelEntity
	Email            string `gorm:"column:email;type:varchar(255) not null;comment:邮箱"`
	PhoneCountryCode string `gorm:"column:phone_country_code;type:varchar(8) not null;comment:手机国家代码"`
	PhoneNumber      string `gorm:"column:phone_number;type:varchar(20) not null;comment:手机号码"`
	Password         string `gorm:"column:password;type:varchar(255) not null;comment:密码"`
	Nickname         string `gorm:"column:nickname;type:varchar(255) not null;comment:昵称"`
	Avatar           string `gorm:"column:avatar;type:varchar(255) not null;comment:头像相对地址"`
	Role             int    `gorm:"column:role;type:int not null;comment:角色：0：普通用户，1：vip"`
	Status           int    `gorm:"column:status;type:tinyint not null;default:0;comment:状态：1：enable，0：disable，-1：deleted"`
}

func (*Entity) TableName() string {
	return cfg.ConfigEntity.TableName
}
