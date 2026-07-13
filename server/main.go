package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"monitor-server/config"
	"monitor-server/internal/dao"
	"monitor-server/internal/handler"
	"monitor-server/internal/middleware"
	"monitor-server/internal/model"
	"monitor-server/internal/scheduler"
	"monitor-server/internal/service"
	"monitor-server/pkg/crypto"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	if err := config.Load("config/config.yaml"); err != nil {
		log.Printf("[Init] 配置加载失败，使用默认: %v", err)
	}
	cfg := config.Get()
	crypto.SetKey(cfg.Crypto.SecretKey)

	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Warn)})
	if err != nil {
		log.Fatalf("[Init] DB失败: %v", err)
	}
	db.Exec("PRAGMA foreign_keys = ON")
	db.AutoMigrate(&model.MonitorTask{}, &model.MonitorRecord{}, &model.ScanRule{}, &model.SysConfig{})
	os.MkdirAll("data", 0755)
	os.MkdirAll("logs", 0755)

	taskDAO := dao.NewTaskDAO(db)
	recordDAO := dao.NewRecordDAO(db)
	configService := service.NewConfigService(db)

	sched := scheduler.New(taskDAO, recordDAO, configService)
	sched.Start()

	taskH := handler.NewTaskHandler(db)
	ruleH := handler.NewRuleHandler(db)
	recordH := handler.NewRecordHandler(db)
	configH := handler.NewConfigHandler(db)
	scanH := handler.NewScanHandler(db)

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()
	r.Use(middleware.CORS(), middleware.Logger())

	api := r.Group("/api")
	tasks := api.Group("/tasks")
	tasks.GET("", taskH.List)
	tasks.POST("", taskH.Create)
	tasks.GET("/:id", taskH.GetByID)
	tasks.PUT("/:id", taskH.Update)
	tasks.DELETE("/:id", taskH.Delete)
	tasks.POST("/:id/execute", scanH.Execute)
	tasks.GET("/:id/rules", ruleH.ListByTask)
	tasks.POST("/:id/rules", ruleH.Create)
	tasks.GET("/:id/records", recordH.ListByTask)

	api.PUT("/rules/:id", ruleH.Update)
	api.DELETE("/rules/:id", ruleH.Delete)
	api.POST("/rules/test", scanH.TestRules)
	api.GET("/records/:id", recordH.GetDetail)
	api.GET("/config/smtp", configH.GetSMTP)
	api.PUT("/config/smtp", configH.SaveSMTP)
	api.POST("/config/smtp/test", configH.TestSMTP)
	api.POST("/cache", scanH.CachePage)

	// SPA fallback
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			c.JSON(404, gin.H{"code": 404, "message": "not found"})
			return
		}
		c.File("web/dist/index.html")
	})

	log.Printf("[Init] 服务启动于 :%d", cfg.Server.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
