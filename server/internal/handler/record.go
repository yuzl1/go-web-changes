package handler

import (
	"monitor-server/internal/model"
	"monitor-server/internal/service"
	"monitor-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecordHandler struct {
	recordService *service.RecordService
}

func NewRecordHandler(db *gorm.DB) *RecordHandler {
	return &RecordHandler{
		recordService: service.NewRecordService(db),
	}
}

// RecordVO 历史列表项（含内容预览和邮件状态）
type RecordVO struct {
	model.MonitorRecord
	ScanPreview string `json:"scan_preview"`
}

// ListByTask 历史列表
func (h *RecordHandler) ListByTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	records, total, err := h.recordService.ListByTaskID(taskID, page, pageSize)
	if err != nil {
		response.InternalError(c, "查询历史记录失败")
		return
	}

	// 添加内容预览（截取前 200 字符）
	list := make([]RecordVO, len(records))
	for i, r := range records {
		preview := r.ScanResult
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		list[i] = RecordVO{MonitorRecord: r, ScanPreview: preview}
	}
	response.Page(c, list, total, page, pageSize)
}

// GetDetail 历史详情
func (h *RecordHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的记录ID")
		return
	}
	detail, err := h.recordService.GetDetail(id)
	if err != nil {
		response.Error(c, 404, "记录不存在")
		return
	}
	response.Success(c, detail)
}
