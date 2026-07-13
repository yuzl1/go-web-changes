package scheduler

import (
	"log"
	"monitor-server/internal/dao"
	"monitor-server/internal/engine"
	"monitor-server/internal/model"
	"monitor-server/internal/notifier"
	"monitor-server/internal/service"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron          *cron.Cron
	taskDAO       *dao.TaskDAO
	recordDAO     *dao.RecordDAO
	configService *service.ConfigService
	running       sync.Map
}

func New(taskDAO *dao.TaskDAO, recordDAO *dao.RecordDAO, configService *service.ConfigService) *Scheduler {
	return &Scheduler{cron: cron.New(cron.WithSeconds()), taskDAO: taskDAO, recordDAO: recordDAO, configService: configService}
}

func (s *Scheduler) Start() error {
	tasks, _ := s.taskDAO.GetAllEnabled()
	for _, t := range tasks {
		s.registerTask(t)
	}
	s.cron.Start()
	log.Printf("[Scheduler] 已启动 %d 个任务", len(tasks))
	return nil
}

func (s *Scheduler) Stop() { <-s.cron.Stop().Done() }

func (s *Scheduler) registerTask(task model.MonitorTask) {
	if cronExpr := model.FreqCron(task.FreqCode); cronExpr != "" {
		taskID := task.ID
		s.cron.AddFunc(cronExpr, func() { s.scanTask(taskID) })
	}
}

func (s *Scheduler) scanTask(taskID int64) {
	if _, exists := s.running.LoadOrStore(taskID, true); exists {
		return
	}
	defer s.running.Delete(taskID)

	task, err := s.taskDAO.GetByID(taskID)
	if err != nil || task.Status != 1 {
		return
	}

	result, scanErr := engine.ExecuteScan(task.TargetURL, task.Rules)
	now := time.Now()
	record := &model.MonitorRecord{TaskID: task.ID, TargetURL: task.TargetURL, ScanTime: now}

	if scanErr != nil {
		record.IsChanged = -1
		record.ErrorMsg = scanErr.Error()
		s.recordDAO.Create(record)
		return
	}

	record.ScanResult = result
	if engine.IsChanged(result, task.LastScanContent) {
		record.IsChanged = 1
		s.taskDAO.UpdateScanResult(task.ID, now, result)
		if task.EmailNotify == 1 {
			s.sendNotification(task, result, now)
			record.EmailSent = 1
		}
	} else {
		record.IsChanged = 0
		s.taskDAO.UpdateScanTime(task.ID, now)
	}
	s.recordDAO.Create(record)
}

func (s *Scheduler) sendNotification(task *model.MonitorTask, content string, t time.Time) {
	cfg, _ := s.configService.GetSMTPConfig()
	if cfg == nil || cfg.Host == "" {
		return
	}
	password, _ := s.configService.GetDecryptedPassword()
	cfg.Password = password
	notifier.SendChangeNotification(cfg, &notifier.Notification{
		Title: "网页变更通知", TaskName: task.Name, TaskURL: task.TargetURL,
		ScanTime: t.Format("2006-01-02 15:04:05"), Content: content,
	})
}
