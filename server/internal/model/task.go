package model

import "time"

type MonitorTask struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"column:name;type:varchar(100);not null" json:"name"`
	TargetURL       string     `gorm:"column:target_url;type:varchar(2048);not null" json:"target_url"`
	FreqCode        int        `gorm:"column:freq_code;not null;default:4" json:"freq_code"`
	Status          int        `gorm:"column:status;not null;default:1" json:"status"`
	EmailNotify     int        `gorm:"column:email_notify;not null;default:1" json:"email_notify"`
	LastScanTime    *time.Time `gorm:"column:last_scan_time" json:"last_scan_time"`
	LastScanContent string     `gorm:"column:last_scan_content;type:text" json:"last_scan_content"`
	Remark          string     `gorm:"column:remark;type:varchar(500)" json:"remark"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	Rules           []ScanRule `gorm:"foreignKey:TaskID" json:"rules,omitempty"`
}

func (MonitorTask) TableName() string { return "monitor_task" }

var FreqCodeMap = map[int]struct{ Name, Cron string }{
	1: {"每5分钟", "*/5 * * * *"},
	2: {"每15分钟", "*/15 * * * *"},
	3: {"每30分钟", "*/30 * * * *"},
	4: {"每1小时", "0 * * * *"},
	5: {"每6小时", "0 */6 * * *"},
	6: {"每12小时", "0 */12 * * *"},
	7: {"每天1次", "0 0 * * *"},
}

func FreqDesc(code int) string {
	if v, ok := FreqCodeMap[code]; ok {
		return v.Name
	}
	return ""
}

func FreqCron(code int) string {
	if v, ok := FreqCodeMap[code]; ok {
		return v.Cron
	}
	return ""
}
