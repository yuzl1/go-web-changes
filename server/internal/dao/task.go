package dao

import (
	"monitor-server/internal/model"

	"gorm.io/gorm"
)

type TaskDAO struct{ db *gorm.DB }

func NewTaskDAO(db *gorm.DB) *TaskDAO { return &TaskDAO{db: db} }

func (d *TaskDAO) List(keyword string, status *int, page, pageSize int) ([]model.MonitorTask, int64, error) {
	var tasks []model.MonitorTask
	var total int64
	query := d.db.Model(&model.MonitorTask{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	query.Count(&total)
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks)
	return tasks, total, query.Error
}

func (d *TaskDAO) GetByID(id int64) (*model.MonitorTask, error) {
	var task model.MonitorTask
	err := d.db.Preload("Rules", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).First(&task, id).Error
	return &task, err
}

func (d *TaskDAO) Create(task *model.MonitorTask) error { return d.db.Create(task).Error }

func (d *TaskDAO) Delete(id int64) error { return d.db.Delete(&model.MonitorTask{}, id).Error }

func (d *TaskDAO) UpdateScanResult(id int64, scanTime interface{}, scanContent string) error {
	return d.db.Model(&model.MonitorTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_scan_time": scanTime, "last_scan_content": scanContent,
	}).Error
}

func (d *TaskDAO) UpdateScanTime(id int64, scanTime interface{}) error {
	return d.db.Model(&model.MonitorTask{}).Where("id = ?", id).Update("last_scan_time", scanTime).Error
}

func (d *TaskDAO) GetAllEnabled() ([]model.MonitorTask, error) {
	var tasks []model.MonitorTask
	err := d.db.Where("status = ?", 1).Preload("Rules", func(db *gorm.DB) *gorm.DB {
		return db.Order("step_order ASC")
	}).Find(&tasks).Error
	return tasks, err
}

func (d *TaskDAO) UpdateCache(id int64, html string) error {
	return d.db.Model(&model.MonitorTask{}).Where("id = ?", id).Update("cached_html", html).Error
}
