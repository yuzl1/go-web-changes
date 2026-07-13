package handler

import (
	"monitor-server/internal/dao"
	"monitor-server/internal/engine"
	"monitor-server/internal/model"
	"monitor-server/internal/notifier"
	"monitor-server/internal/service"
	"monitor-server/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScanHandler struct {
	taskDAO       *dao.TaskDAO
	recordDAO     *dao.RecordDAO
	configService *service.ConfigService
}

func NewScanHandler(db *gorm.DB) *ScanHandler {
	return &ScanHandler{
		taskDAO:       dao.NewTaskDAO(db),
		recordDAO:     dao.NewRecordDAO(db),
		configService: service.NewConfigService(db),
	}
}

func (h *ScanHandler) Execute(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	task, err := h.taskDAO.GetByID(id)
	if err != nil { response.Error(c, 404, "任务不存在"); return }
	result, scanErr := engine.ExecuteScan(task.TargetURL, task.Rules)
	now := time.Now()
	record := &model.MonitorRecord{TaskID: task.ID, TargetURL: task.TargetURL, ScanTime: now}
	if scanErr != nil {
		record.IsChanged = -1; record.ErrorMsg = scanErr.Error()
		h.recordDAO.Create(record)
		response.Success(c, map[string]interface{}{"record_id": record.ID, "is_changed": -1, "error": scanErr.Error()})
		return
	}
	record.ScanResult = result
	if engine.IsChanged(result, task.LastScanContent) {
		record.IsChanged = 1; h.taskDAO.UpdateScanResult(task.ID, now, result)
		if task.EmailNotify == 1 { go h.sendChangeNotification(task, result, now); record.EmailSent = 1 }
	} else { record.IsChanged = 0; h.taskDAO.UpdateScanTime(task.ID, now) }
	h.recordDAO.Create(record)
	response.Success(c, map[string]interface{}{"record_id": record.ID, "is_changed": record.IsChanged, "scan_result": result, "scan_time": now})
}

func (h *ScanHandler) TestRules(c *gin.Context) {
	var req struct{ TargetURL string `json:"target_url"`; Rules []model.ScanRule `json:"rules"` }
	if err := c.ShouldBindJSON(&req); err != nil || req.TargetURL == "" || len(req.Rules) == 0 { response.BadRequest(c, "参数错误"); return }
	steps, final, err := engine.ExecuteTest(req.TargetURL, req.Rules)
	msg := "测试完成"
	if err != nil { msg = "存在失败步骤" }
	response.Success(c, map[string]interface{}{"steps": steps, "final_result": final, "message": msg})
}

func (h *ScanHandler) CachePage(c *gin.Context) {
	var req struct{ TargetURL string `json:"target_url"` }
	c.ShouldBindJSON(&req)
	if req.TargetURL == "" { response.BadRequest(c, "URL不能为空"); return }
	html, err := engine.FetchPageHTML(req.TargetURL)
	if err != nil { response.Error(c, 500, "获取失败: "+err.Error()); return }
	response.Success(c, map[string]interface{}{"html": html, "html_length": len(html)})
}

func (h *ScanHandler) sendChangeNotification(task *model.MonitorTask, content string, t time.Time) {
	cfg, _ := h.configService.GetSMTPConfig()
	pwd, _ := h.configService.GetDecryptedPassword(); cfg.Password = pwd
	notifier.SendChangeNotification(cfg, &notifier.Notification{Title: "网页变更通知", TaskName: task.Name, TaskURL: task.TargetURL, ScanTime: t.Format("2006-01-02 15:04:05"), Content: content})
}
