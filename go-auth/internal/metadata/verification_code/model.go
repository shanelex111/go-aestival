package verification_code

import "go-auth/internal/base"

const (
	StatusPending = "pending"
	StatusUsed    = "used"
	StatusExpired = "expired"

	SceneSignin        = "signin"
	SceneResetPassword = "reset_password"
)

type VerificationCodeEntity struct {
	base.BaseModelEntity
	Scene       string `gorm:"column:scene;type:varchar(100) not null;comment:场景: signin | reset_password"`
	Type        string `gorm:"column:type;type:varchar(10) not null;comment:类型，email | phone"`
	Target      string `gorm:"column:target;type:varchar(100) not null;comment:目标，邮箱或手机号"`
	CountryCode string `gorm:"column:country_code;type:varchar(8) not null;comment:国家代码"`
	Code        string `gorm:"column:code;type:varchar(10) not null;comment:验证码"`
	Status      string `gorm:"column:status;type:varchar(20) not null;comment:状态：pending | used | expired"`
	ExpiredAt   int64  `gorm:"column:expired_at;type:bigint unsigned not null;comment:过期时间，毫秒时间戳"`
	Platform    string `gorm:"column:platform;type:varchar(20) not null;comment:第三方：aliyun | tencent"`
	TemplateID  string `gorm:"column:template_id;type:varchar(255) not null;comment:模板id"`
	Content     string `gorm:"column:content;type:text not null;comment:内容"`
}

func (VerificationCodeEntity) TableName() string {
	return cfg.ConfigEntity.TableName
}
