package dao

import (
	"monitor-server/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConfigDAO struct {
	db *gorm.DB
}

func NewConfigDAO(db *gorm.DB) *ConfigDAO {
	return &ConfigDAO{db: db}
}

// GetByKey 按 key 查询配置
func (d *ConfigDAO) GetByKey(key string) (*model.SysConfig, error) {
	var cfg model.SysConfig
	if err := d.db.Where("config_key = ?", key).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetByKeys 按 key 批量查询配置
func (d *ConfigDAO) GetByKeys(keys []string) (map[string]string, error) {
	var configs []model.SysConfig
	if err := d.db.Where("config_key IN ?", keys).Find(&configs).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, c := range configs {
		result[c.ConfigKey] = c.ConfigValue
	}
	return result, nil
}

// Upsert 插入或更新配置
func (d *ConfigDAO) Upsert(key, value string) error {
	cfg := model.SysConfig{
		ConfigKey:   key,
		ConfigValue: value,
	}
	return d.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "config_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"config_value", "updated_at"}),
	}).Create(&cfg).Error
}
