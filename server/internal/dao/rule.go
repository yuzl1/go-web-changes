package dao

import (
	"monitor-server/internal/model"
	"gorm.io/gorm"
)

type RuleDAO struct{ db *gorm.DB }
func NewRuleDAO(db *gorm.DB) *RuleDAO { return &RuleDAO{db: db} }

func (d *RuleDAO) GetByTaskID(tid int64) ([]model.ScanRule, error) {
	var rules []model.ScanRule
	return rules, d.db.Where("task_id=?", tid).Order("step_order ASC").Find(&rules).Error
}
func (d *RuleDAO) BatchCreate(rules []model.ScanRule) error { if len(rules)==0{return nil}; return d.db.Create(&rules).Error }
func (d *RuleDAO) Update(rule *model.ScanRule) error {
	return d.db.Model(&model.ScanRule{}).Where("id=?", rule.ID).Updates(map[string]interface{}{
		"step_order": rule.StepOrder, "rule_content": rule.RuleContent, "rule_mode": rule.RuleMode,
	}).Error
}
func (d *RuleDAO) Delete(id int64) error { return d.db.Delete(&model.ScanRule{}, id).Error }
func (d *RuleDAO) GetByID(id int64) (*model.ScanRule, error) {
	var r model.ScanRule; return &r, d.db.First(&r, id).Error
}
