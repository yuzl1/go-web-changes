package handler

import (
	"monitor-server/internal/model"
	"monitor-server/internal/service"
	"monitor-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct{ taskService *service.TaskService }
func NewTaskHandler(db *gorm.DB) *TaskHandler { return &TaskHandler{service.NewTaskService(db)} }

func (h *TaskHandler) List(c *gin.Context) {
	kw := c.Query("keyword"); page, _ := strconv.Atoi(c.DefaultQuery("page", "1")); ps, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	var status *int
	if s := c.Query("status"); s != "" { v, _ := strconv.Atoi(s); status = &v }
	tasks, total, _ := h.taskService.List(kw, status, page, ps)
	type VO struct{ model.MonitorTask; FreqDesc string `json:"freq_desc"` }
	list := make([]VO, len(tasks))
	for i, t := range tasks { list[i] = VO{MonitorTask: t, FreqDesc: model.FreqDesc(t.FreqCode)} }
	response.Page(c, list, total, page, ps)
}
func (h *TaskHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	task, err := h.taskService.GetByID(id)
	if err != nil { response.Error(c, 404, "不存在"); return }
	response.Success(c, task)
}
func (h *TaskHandler) Create(c *gin.Context) {
	var req struct {
		Name, TargetURL, Remark string; FreqCode, Status, EmailNotify int; Rules []model.ScanRule
	}
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
	task := &model.MonitorTask{Name: req.Name, TargetURL: req.TargetURL, FreqCode: req.FreqCode, Status: req.Status, EmailNotify: req.EmailNotify, Remark: req.Remark}
	if err := h.taskService.Create(task, req.Rules); err != nil { response.BadRequest(c, err.Error()); return }
	response.SuccessMsg(c, "创建成功")
}
func (h *TaskHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Name, TargetURL, Remark string; FreqCode, Status, EmailNotify int; Rules []model.ScanRule
	}
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
	task := &model.MonitorTask{ID: id, Name: req.Name, TargetURL: req.TargetURL, FreqCode: req.FreqCode, Status: req.Status, EmailNotify: req.EmailNotify, Remark: req.Remark}
	if err := h.taskService.Update(task, req.Rules); err != nil { response.BadRequest(c, err.Error()); return }
	response.SuccessMsg(c, "更新成功")
}
func (h *TaskHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.taskService.Delete(id)
	response.SuccessMsg(c, "删除成功")
}
