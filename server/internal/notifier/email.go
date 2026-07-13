package notifier

import (
	"errors"
	"log"
	"monitor-server/internal/model"
	"strings"

	gomail "gopkg.in/gomail.v2"
)

type Notification struct {
	Title    string
	TaskName string
	TaskURL  string
	ScanTime string
	Content  string
}

// SendTestEmail 发送测试邮件
func SendTestEmail(cfg *model.SMTPConfig) error {
	if cfg.Host == "" {
		return errors.New("SMTP服务器未配置")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(cfg.SenderEmail, cfg.SenderName))
	m.SetHeader("To", strings.Split(cfg.ReceiverEmails, ",")...)
	m.SetHeader("Subject", "[网页变更监听] 测试邮件")
	m.SetBody("text/html", `<h3>测试邮件</h3><p>这是一封来自网页变更监听系统的测试邮件。</p><p>如果您收到此邮件，说明 SMTP 配置正确。</p>`)

	return send(cfg, m)
}

// SendChangeNotification 发送变更通知
func SendChangeNotification(cfg *model.SMTPConfig, n *Notification) error {
	if cfg == nil || cfg.Host == "" {
		return errors.New("SMTP服务器未配置")
	}
	if n == nil {
		return errors.New("通知内容为空")
	}

	receivers := strings.Split(cfg.ReceiverEmails, ",")
	if len(receivers) == 0 || cfg.ReceiverEmails == "" {
		return errors.New("收件人邮箱未配置")
	}

	// 截取内容前500字符
	content := n.Content
	if len(content) > 500 {
		content = content[:500] + "..."
	}
	// HTML 转义
	content = strings.ReplaceAll(content, "<", "&lt;")
	content = strings.ReplaceAll(content, ">", "&gt;")

	body := buildEmailBody(n.TaskName, n.TaskURL, n.ScanTime, content)

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(cfg.SenderEmail, cfg.SenderName))
	m.SetHeader("To", receivers...)
	m.SetHeader("Subject", "[网页变更通知] "+n.TaskName+" — 检测到内容变更")
	m.SetBody("text/html", body)

	return send(cfg, m)
}

// send 发送邮件
func send(cfg *model.SMTPConfig, m *gomail.Message) error {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	switch cfg.Encryption {
	case "PLAIN":
		// 无 TLS
		d.TLSConfig = nil
		d.SSL = false
	case "TLS":
		if cfg.Port == 465 {
			// 隐式 TLS (SMTPS)
			d.SSL = true
		}
		// port 587 时 gomail 默认使用 STARTTLS
	default:
		d.TLSConfig = nil
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("[Email] 发送失败: %v", err)
		return err
	}
	log.Printf("[Email] 发送成功: %s", cfg.ReceiverEmails)
	return nil
}

// buildEmailBody 构建邮件 HTML 正文
func buildEmailBody(taskName, taskURL, scanTime, content string) string {
	return `<!DOCTYPE html>
<html>
<body>
<h2>网页变更通知</h2>
<p>您监控的网页已检测到内容变更：</p>
<table border="1" cellpadding="8" cellspacing="0" style="border-collapse:collapse;">
  <tr><td><b>任务名称</b></td><td>` + taskName + `</td></tr>
  <tr><td><b>目标URL</b></td><td><a href="` + taskURL + `">` + taskURL + `</a></td></tr>
  <tr><td><b>扫描时间</b></td><td>` + scanTime + `</td></tr>
  <tr><td><b>变更内容</b></td><td><pre style="max-width:600px;overflow:auto;">` + content + `</pre></td></tr>
</table>
<p>请登录系统查看完整详情。</p>
</body>
</html>`
}
