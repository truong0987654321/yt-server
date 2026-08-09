package router

import (
	"yt-clone-server/internal/handler"
	"yt-clone-server/internal/middleware"
	"yt-clone-server/internal/service"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
	JWTService  *service.JWTService
}

func Setup(deps Dependencies) *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())

	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.GET("/google/login", deps.AuthHandler.GoogleLogin)         // web
		auth.GET("/google/callback", deps.AuthHandler.GoogleCallback)   // web
		auth.POST("/google/mobile", deps.AuthHandler.GoogleMobileLogin) // mobile
		auth.POST("/refresh", deps.AuthHandler.RefreshToken)            // dùng chung
		auth.POST("/logout", deps.AuthHandler.Logout)
	}

	users := api.Group("/users")
	users.Use(middleware.RequireAuth(deps.JWTService)) // mọi route trong nhóm này bắt buộc có JWT
	{
		users.GET("/me", deps.UserHandler.Me)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // production: đổi thành domain Next.js cụ thể
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
