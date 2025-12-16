package api

import (
	"github.com/gin-gonic/gin"

	"tidalcore-backend/middleware"
)

func SetupRouter(mode string) *gin.Engine {
	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	userHandler := NewUserHandler()
	checkinHandler := NewCheckinHandler()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 公开接口
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// 公开数据
		v1.GET("/leaderboard", userHandler.GetLeaderboard)
		v1.GET("/heatmap/global", checkinHandler.GetGlobalHeatmap)

		// 需要认证的接口
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// 用户相关
			protected.GET("/user/profile", userHandler.GetProfile)

			// 打卡相关
			protected.POST("/checkin", checkinHandler.Checkin)
			protected.GET("/checkin/history", checkinHandler.GetHistory)
			protected.GET("/checkin/heatmap", checkinHandler.GetHeatmap)
		}
	}

	return r
}
