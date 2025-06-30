package aestival_member

import "go-auth/internal/base"

type Entity struct {
	base.BaseModelEntity
	AccountID uint   `gorm:"column:account_id;type:int unsigned not null;comment:账户id"`
	Nickname  string `gorm:"column:nickname;type:varchar(255) not null;comment:昵称"`
	Avatar    string `gorm:"column:avatar;type:varchar(255) not null;comment:头像相对地址"`
	Role      int    `gorm:"column:role;type:int not null;comment:角色：0：普通用户，1：vip"`
}

func (*Entity) TableName() string {
	return cfg.ConfigEntity.TableName
}
