package main

import (
	"log"
	"yt-clone-server/internal/config"
	"yt-clone-server/internal/handler"
	"yt-clone-server/internal/repository/postgres"
	"yt-clone-server/internal/router"
	"yt-clone-server/internal/service"
)

func main() {
	log.Println("Starting the Go Backend Server...")

	// 1. Load config
	cfg := config.Load()

	// 2. Kết nối DB
	db, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// 3. Khởi tạo repository (tầng data access)
	userRepo := postgres.NewUserRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	channelRepo := postgres.NewChannelRepository(db)

	// 4. Khởi tạo service (tầng business logic)
	jwtService := service.NewJWTService(cfg)
	authService := service.NewAuthService(cfg, userRepo, refreshTokenRepo, jwtService)
	categoryService := service.NewCategoryService(categoryRepo)
	channelService := service.NewChannelService(channelRepo, userRepo)

	// 5. Khởi tạo handler (tầng HTTP)
	authHandler := handler.NewAuthHandler(cfg, authService)
	userHandler := handler.NewUserHandler(userRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	channelHandler := handler.NewChannelHandler(channelService)

	// 6. Setup router
	r := router.Setup(router.Dependencies{
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		CategoryHandler: categoryHandler,
		ChannelHandler:  channelHandler,
		JWTService:      jwtService,
	})
	log.Printf("Server is running on port %s (env: %s)", cfg.Port, cfg.Env)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
