package model

import "time"

// SysConfig 系统配置表
type SysConfig struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string    `gorm:"column:config_key;type:varchar(100);not null;uniqueIndex" json:"config_key"`
	ConfigValue string    `gorm:"column:config_value;type:text;not null" json:"config_value"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}

// 预设 SMTP 配置项 key
const (
	ConfigKeySMTPHost           = "smtp_host"
	ConfigKeySMTPPort           = "smtp_port"
	ConfigKeySMTPEncryption     = "smtp_encryption"
	ConfigKeySMTPSenderEmail    = "smtp_sender_email"
	ConfigKeySMTPSenderName     = "smtp_sender_name"
	ConfigKeySMTPUsername       = "smtp_username"
	ConfigKeySMTPPassword       = "smtp_password"
	ConfigKeySMTPReceiverEmails = "smtp_receiver_emails"
)

var SMTPConfigKeys = []string{
	ConfigKeySMTPHost,
	ConfigKeySMTPPort,
	ConfigKeySMTPEncryption,
	ConfigKeySMTPSenderEmail,
	ConfigKeySMTPSenderName,
	ConfigKeySMTPUsername,
	ConfigKeySMTPPassword,
	ConfigKeySMTPReceiverEmails,
}

// SMTPConfig 邮件配置结构
type SMTPConfig struct {
	Host           string `json:"smtp_host"`
	Port           int    `json:"smtp_port"`
	Encryption     string `json:"smtp_encryption"` // PLAIN / TLS
	SenderEmail    string `json:"smtp_sender_email"`
	SenderName     string `json:"smtp_sender_name"`
	Username       string `json:"smtp_username"`
	Password       string `json:"smtp_password"`
	ReceiverEmails string `json:"smtp_receiver_emails"`
}
