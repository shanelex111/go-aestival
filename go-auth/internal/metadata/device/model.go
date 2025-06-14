package device

import "go-auth/internal/base"

type DeviceEntity struct {
	base.BaseModelEntity
	AccountID                uint   `gorm:"column:account_id;type:int unsigned not null;comment:账户id"`
	DeviceID                 string `gorm:"column:device_id;type:varchar(255) not null;comment:设备id"`
	DeviceType               string `gorm:"column:device_type;type:varchar(10) not null;comment:设备类型"`
	DeviceModel              string `gorm:"column:device_model;type:varchar(255) not null;comment:设备型号"`
	AppVersion               string `gorm:"column:app_version;type:int unsigned not null;comment:app版本"`
	SigninTimes              int    `gorm:"column:signin_times;type:int unsigned not null;default:0;comment:该设备总登录次数"`
	CreatedIP                string `gorm:"column:created_ip;type:varchar(45) not null;comment:首次使用该设备登录时IP"`
	CreatedIPCountryCode     string `gorm:"column:created_ip_country_code;type:varchar(8) not null;comment:首次使用该设备登录时IP所属国家"`
	CreatedIPSubdivisionCode string `gorm:"column:created_ip_subdivision_code;type:varchar(8) not null;comment:首次使用该设备登录时IP所属地区"`
	UpdatedIP                string `gorm:"column:updated_ip;type:varchar(45) not null;comment:最近使用该设备登录时IP"`
	UpdatedIPCountryCode     string `gorm:"column:updated_ip_country_code;type:varchar(8) not null;comment:最近使用该设备登录时IP所属国家"`
	UpdatedIPSubdivisionCode string `gorm:"column:updated_ip_subdivision_code;type:varchar(8) not null;comment:最近使用该设备登录时IP所属地区"`
}

func (DeviceEntity) TableName() string {
	return cfg.ConfigEntity.TableName
}
