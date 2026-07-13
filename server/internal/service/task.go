package service

import (
	"errors"
	"monitor-server/internal/dao"
	"monitor-server/internal/model"
	"gorm.io/gorm"
)

type TaskService struct {
	taskDAO *dao.TaskDAO
	ruleDAO *dao.RuleDAO
	db      *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{taskDAO: dao.NewTaskDAO(db), ruleDAO: dao.NewRuleDAO(db), db: db}
}

func (s *TaskService) List(kw string, status *int, page, ps int) ([]model.MonitorTask, int64, error) {
	return s.taskDAO.List(kw, status, page, ps)
}
func (s *TaskService) GetByID(id int64) (*model.MonitorTask, error) { return s.taskDAO.GetByID(id) }

func (s *TaskService) Create(task *model.MonitorTask, rules []model.ScanRule) error {
	if task.Name == "" { return errors.New("名称不能为空") }
	if task.TargetURL == "" || (len(task.TargetURL)<8 || (task.TargetURL[:7]!="http://"&&task.TargetURL[:8]!="https://")) { return errors.New("URL非法") }
	if len(rules) == 0 { return errors.New("至少一条规则") }
	for _, r := range rules {
		if r.RuleContent == "" { return errors.New("脚本不能为空") }
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(task).Error; err != nil { return err }
		for i := range rules { rules[i].TaskID = task.ID }
		if len(rules) > 0 { return tx.Create(&rules).Error }
		return nil
	})
}

func (s *TaskService) Update(task *model.MonitorTask, rules []model.ScanRule) error {
	if task.ID <= 0 { return errors.New("ID为空") }
	if task.Name == "" { return errors.New("名称不能为空") }
	if len(rules) == 0 { return errors.New("至少一条规则") }
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&model.MonitorTask{}).Where("id=?", task.ID).Updates(map[string]interface{}{
			"name": task.Name, "target_url": task.TargetURL, "freq_code": task.FreqCode,
			"status": task.Status, "email_notify": task.EmailNotify, "remark": task.Remark,
		})
		tx.Where("task_id=?", task.ID).Delete(&model.ScanRule{})
		for i := range rules { rules[i].ID = 0; rules[i].TaskID = task.ID }
		if len(rules) > 0 { return tx.Create(&rules).Error }
		return nil
	})
}

func (s *TaskService) Delete(id int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("task_id=?", id).Delete(&model.ScanRule{})
		tx.Where("task_id=?", id).Delete(&model.MonitorRecord{})
		return tx.Delete(&model.MonitorTask{}, id).Error
	})
}
