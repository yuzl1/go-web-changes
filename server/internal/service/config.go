package service

import (
	"monitor-server/internal/dao"
	"monitor-server/internal/model"
	"monitor-server/pkg/crypto"
	"strconv"

	"gorm.io/gorm"
)

type ConfigService struct {
	configDAO *dao.ConfigDAO
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{
		configDAO: dao.NewConfigDAO(db),
	}
}

// GetSMTPConfig 获取 SMTP 配置
func (s *ConfigService) GetSMTPConfig() (*model.SMTPConfig, error) {
	kv, err := s.configDAO.GetByKeys(model.SMTPConfigKeys)
	if err != nil {
		return nil, err
	}

	port, _ := strconv.Atoi(kv[model.ConfigKeySMTPPort])

	return &model.SMTPConfig{
		Host:           kv[model.ConfigKeySMTPHost],
		Port:           port,
		Encryption:     kv[model.ConfigKeySMTPEncryption],
		SenderEmail:    kv[model.ConfigKeySMTPSenderEmail],
		SenderName:     kv[model.ConfigKeySMTPSenderName],
		Username:       kv[model.ConfigKeySMTPUsername],
		Password:       "", // 不返回密码
		ReceiverEmails: kv[model.ConfigKeySMTPReceiverEmails],
	}, nil
}

// SaveSMTPConfig 保存 SMTP 配置
func (s *ConfigService) SaveSMTPConfig(cfg *model.SMTPConfig) error {
	upsert := func(key, value string) error {
		return s.configDAO.Upsert(key, value)
	}

	if err := upsert(model.ConfigKeySMTPHost, cfg.Host); err != nil {
		return err
	}
	if err := upsert(model.ConfigKeySMTPPort, strconv.Itoa(cfg.Port)); err != nil {
		return err
	}
	if err := upsert(model.ConfigKeySMTPEncryption, cfg.Encryption); err != nil {
		return err
	}
	if err := upsert(model.ConfigKeySMTPSenderEmail, cfg.SenderEmail); err != nil {
		return err
	}
	if err := upsert(model.ConfigKeySMTPSenderName, cfg.SenderName); err != nil {
		return err
	}
	if err := upsert(model.ConfigKeySMTPUsername, cfg.Username); err != nil {
		return err
	}
	if err := upsert(model.ConfigKeySMTPReceiverEmails, cfg.ReceiverEmails); err != nil {
		return err
	}

	// 密码非空时才更新
	if cfg.Password != "" {
		encrypted, err := crypto.Encrypt(cfg.Password)
		if err != nil {
			return err
		}
		if err := upsert(model.ConfigKeySMTPPassword, encrypted); err != nil {
			return err
		}
	}

	return nil
}

// GetDecryptedPassword 获取解密后的 SMTP 密码
func (s *ConfigService) GetDecryptedPassword() (string, error) {
	cfg, err := s.configDAO.GetByKey(model.ConfigKeySMTPPassword)
	if err != nil {
		return "", err
	}
	if cfg.ConfigValue == "" {
		return "", nil
	}
	return crypto.Decrypt(cfg.ConfigValue)
}
