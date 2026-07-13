package dao

import (
	"monitor-server/internal/model"

	"gorm.io/gorm"
)

type RecordDAO struct {
	db *gorm.DB
}

func NewRecordDAO(db *gorm.DB) *RecordDAO {
	return &RecordDAO{db: db}
}

// ListByTaskID 按任务ID分页查询历史记录（倒序）
func (d *RecordDAO) ListByTaskID(taskID int64, page, pageSize int) ([]model.MonitorRecord, int64, error) {
	var records []model.MonitorRecord
	var total int64

	query := d.db.Model(&model.MonitorRecord{}).Where("task_id = ?", taskID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("scan_time DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// GetByID 按ID查询记录详情
func (d *RecordDAO) GetByID(id int64) (*model.MonitorRecord, error) {
	var record model.MonitorRecord
	if err := d.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// GetPrevRecord 获取同一任务的上一条扫描记录（用于对比）
func (d *RecordDAO) GetPrevRecord(taskID int64, currentID int64) (*model.MonitorRecord, error) {
	var record model.MonitorRecord
	err := d.db.Where("task_id = ? AND id < ?", taskID, currentID).
		Order("id DESC").First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// Create 创建扫描记录
func (d *RecordDAO) Create(record *model.MonitorRecord) error {
	return d.db.Create(record).Error
}
