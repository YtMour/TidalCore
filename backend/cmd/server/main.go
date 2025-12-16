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
	configPath := flag.String("config", "", "path to config file (optional, can use env vars)")
	flag.Parse()

	// 加载配置 (支持配置文件或环境变量)
	config.Load(*configPath)
	cfg := config.Get()

	log.Printf("Starting TidalCore Backend...")
	log.Printf("Database: %s@%s:%s/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer database.Close()

	// 自动迁移数据库表
	if err := autoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Printf("Database migration completed")

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
