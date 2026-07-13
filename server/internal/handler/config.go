package handler

import (
	"monitor-server/internal/model"
	"monitor-server/internal/notifier"
	"monitor-server/internal/service"
	"monitor-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	configService *service.ConfigService
}

func NewConfigHandler(db *gorm.DB) *ConfigHandler {
	return &ConfigHandler{
		configService: service.NewConfigService(db),
	}
}

// GetSMTP 获取 SMTP 配置
func (h *ConfigHandler) GetSMTP(c *gin.Context) {
	cfg, err := h.configService.GetSMTPConfig()
	if err != nil {
		response.InternalError(c, "获取配置失败")
		return
	}
	response.Success(c, cfg)
}

// SaveSMTP 保存 SMTP 配置
func (h *ConfigHandler) SaveSMTP(c *gin.Context) {
	var cfg model.SMTPConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		response.BadRequest(c, "参数解析失败")
		return
	}
	if err := h.configService.SaveSMTPConfig(&cfg); err != nil {
		response.InternalError(c, "保存配置失败: "+err.Error())
		return
	}
	response.SuccessMsg(c, "保存成功")
}

// TestSMTP 测试发送邮件
func (h *ConfigHandler) TestSMTP(c *gin.Context) {
	cfg, err := h.configService.GetSMTPConfig()
	if err != nil {
		response.InternalError(c, "获取SMTP配置失败")
		return
	}
	// 获取解密后的密码
	password, err := h.configService.GetDecryptedPassword()
	if err != nil {
		password = ""
	}
	cfg.Password = password

	if err := notifier.SendTestEmail(cfg); err != nil {
		response.Error(c, 500, "发送失败: "+err.Error())
		return
	}
	response.SuccessMsg(c, "测试邮件已发送")
}
