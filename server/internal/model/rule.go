package model

import "time"

type ScanRule struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID      int64     `gorm:"column:task_id;not null;index" json:"task_id"`
	StepOrder   int       `gorm:"column:step_order;not null" json:"step_order"`
	RuleContent string    `gorm:"column:rule_content;type:text;not null" json:"rule_content"` // jQuery脚本
	RuleMode    int       `gorm:"column:rule_mode;not null;default:1" json:"rule_mode"`       // 1=必须成功 2=失败跳过
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (ScanRule) TableName() string { return "scan_rule" }

var RuleModeMap = map[int]string{1: "必须成功", 2: "失败跳过"}
