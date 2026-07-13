package service

import (
	"monitor-server/internal/dao"
	"monitor-server/internal/model"

	"gorm.io/gorm"
)

type RecordService struct {
	recordDAO *dao.RecordDAO
}

func NewRecordService(db *gorm.DB) *RecordService {
	return &RecordService{
		recordDAO: dao.NewRecordDAO(db),
	}
}

// ListByTaskID 历史列表
func (s *RecordService) ListByTaskID(taskID int64, page, pageSize int) ([]model.MonitorRecord, int64, error) {
	return s.recordDAO.ListByTaskID(taskID, page, pageSize)
}

// GetDetail 历史详情（含上一条记录用于对比）
func (s *RecordService) GetDetail(id int64) (map[string]interface{}, error) {
	record, err := s.recordDAO.GetByID(id)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"id":          record.ID,
		"task_id":     record.TaskID,
		"target_url":  record.TargetURL,
		"scan_time":   record.ScanTime,
		"scan_result": record.ScanResult,
		"is_changed":  record.IsChanged,
		"error_msg":   record.ErrorMsg,
	}

	// 获取上一条记录
	prev, err := s.recordDAO.GetPrevRecord(record.TaskID, record.ID)
	if err == nil {
		result["prev_scan_result"] = prev.ScanResult
	} else {
		result["prev_scan_result"] = ""
	}

	return result, nil
}
