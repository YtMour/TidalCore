package main

import (
	"flag"
	"fmt"
	"log"

	"tidalcore-backend/api"
	"tidalcore-backend/config"
	"tidalcore-backend/internal/model"
	"tidalcore-backend/pkg/database"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	// 加载配置
	if err := config.Load(*configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg := config.Get()

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer database.Close()

	// 自动迁移数据库表
	if err := autoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 启动服务器
	r := api.SetupRouter(cfg.Server.Mode)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func autoMigrate() error {
	db := database.Get()
	return db.AutoMigrate(
		&model.User{},
		&model.Checkin{},
	)
}
