package model

import "time"

// MonitorRecord 扫描记录子表
type MonitorRecord struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID     int64     `gorm:"column:task_id;not null;index" json:"task_id"`
	TargetURL  string    `gorm:"column:target_url;type:varchar(2048);not null" json:"target_url"`
	ScanTime   time.Time `gorm:"column:scan_time;autoCreateTime" json:"scan_time"`
	ScanResult string    `gorm:"column:scan_result;type:text" json:"scan_result"`
	IsChanged  int       `gorm:"column:is_changed;not null;default:0" json:"is_changed"`   // 1=有变更 0=无变更 -1=失败
	ErrorMsg   string    `gorm:"column:error_msg;type:varchar(1000)" json:"error_msg"`
	EmailSent  int       `gorm:"column:email_sent;not null;default:0" json:"email_sent"`   // 1=已发送 0=未发送 -1=发送失败
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (MonitorRecord) TableName() string {
	return "monitor_record"
}
