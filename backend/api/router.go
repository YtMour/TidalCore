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
	visitHandler := NewVisitHandler()

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

		// 访问统计 (公开)
		v1.POST("/visit", visitHandler.RecordVisit)
		v1.GET("/visit/stats", visitHandler.GetStats)
		v1.GET("/visit/heatmap", visitHandler.GetHeatmap)

		// 需要认证的接口
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// 用户相关
			protected.GET("/user/profile", userHandler.GetProfile)
			protected.PUT("/user/profile", userHandler.UpdateProfile)
			protected.PUT("/user/username", userHandler.UpdateUsername)
			protected.PUT("/user/password", userHandler.UpdatePassword)

			// 打卡相关
			protected.POST("/checkin", checkinHandler.Checkin)
			protected.GET("/checkin/history", checkinHandler.GetHistory)
			protected.GET("/checkin/heatmap", checkinHandler.GetHeatmap)
		}

		// 管理员接口
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuthWithAdmin())
		admin.Use(middleware.AdminAuth())
		{
			admin.GET("/users", userHandler.GetAllUsers)
			admin.DELETE("/users/:id", userHandler.DeleteUser)
			admin.PUT("/users/:id/admin", userHandler.SetUserAdmin)
			admin.PUT("/users/:id/stats", userHandler.UpdateUserStats)
		}
	}

	return r
}
