package base

type BaseModelEntity struct {
	ID        uint  `gorm:"column:id;type:int unsigned not null AUTO_INCREMENT;primaryKey;comment:主键id"`
	CreatedAt int64 `gorm:"column:created_at;autoCreateTime:milli;type:bigint unsigned not null;comment:创建时间，毫秒时间戳"`
	UpdatedAt int64 `gorm:"column:updated_at;autoUpdateTime:milli;type:bigint unsigned not null;comment:更新时间，毫秒时间戳"`
}
