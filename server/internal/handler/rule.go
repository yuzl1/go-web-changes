package handler

import (
	"monitor-server/internal/dao"
	"monitor-server/internal/model"
	"monitor-server/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RuleHandler struct{ ruleDAO *dao.RuleDAO }
func NewRuleHandler(db *gorm.DB) *RuleHandler { return &RuleHandler{dao.NewRuleDAO(db)} }

func (h *RuleHandler) ListByTask(c *gin.Context) {
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	rules, _ := h.ruleDAO.GetByTaskID(tid)
	if rules == nil { rules = []model.ScanRule{} }
	response.Success(c, rules)
}
func (h *RuleHandler) Create(c *gin.Context) {
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct{ StepOrder, RuleMode int; RuleContent string }
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
	r := &model.ScanRule{TaskID: tid, StepOrder: req.StepOrder, RuleContent: req.RuleContent, RuleMode: req.RuleMode}
	if r.RuleMode == 0 { r.RuleMode = 1 }
	h.ruleDAO.BatchCreate([]model.ScanRule{*r})
	response.SuccessMsg(c, "添加成功")
}
func (h *RuleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct{ StepOrder, RuleMode int; RuleContent string }
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
	r := &model.ScanRule{ID: id, StepOrder: req.StepOrder, RuleContent: req.RuleContent, RuleMode: req.RuleMode}
	h.ruleDAO.Update(r)
	response.SuccessMsg(c, "更新成功")
}
func (h *RuleHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.ruleDAO.Delete(id)
	response.SuccessMsg(c, "删除成功")
}
