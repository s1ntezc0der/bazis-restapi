package run

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"mkk_bazis/internal/config"
	auth_controller "mkk_bazis/internal/services/auth/controller"
	auth_repo "mkk_bazis/internal/services/auth/repository"
	auth_usecase "mkk_bazis/internal/services/auth/usecase"
	teams_controller "mkk_bazis/internal/services/teams/controller"
	teams_repo "mkk_bazis/internal/services/teams/repository"
	teams_usecase "mkk_bazis/internal/services/teams/usecase"
	tasks_controller "mkk_bazis/internal/services/tasks/controller"
	tasks_repo "mkk_bazis/internal/services/tasks/repository"
	tasks_usecase "mkk_bazis/internal/services/tasks/usecase"
	"mkk_bazis/pkg/cache"
	"mkk_bazis/pkg/db"
	"mkk_bazis/pkg/jwt"
	"mkk_bazis/pkg/logger"
	"mkk_bazis/pkg/middleware"
)

func Run() error {
	// 1. Загружаем конфиг
	cfg := config.Load()

	// 2. Инициализируем логгер
	logger := logger.New()

	// 3. Подключаемся к MySQL
	dbConn, err := db.NewMySQLDB(db.Config{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Name:     cfg.DB.Name,
	})
	if err != nil {
		return err
	}
	defer dbConn.Close()

	// 4. Инициализируем JWT
	jwtConfig := jwt.NewJWT(cfg.JWT.Secret, time.Duration(cfg.JWT.Expire)*time.Hour)

	// 5. Auth
	authRepo := auth_repo.NewAuthRepository(dbConn)
	authService := auth_usecase.NewAuthService(authRepo, jwtConfig)
	authHandler := auth_controller.NewAuthHandler(authService)

	// 6. Teams
	teamRepo := teams_repo.NewTeamRepository(dbConn)
	teamService := teams_usecase.NewTeamService(teamRepo, authRepo)
	teamHandler := teams_controller.NewTeamHandler(teamService)

	// 7. Tasks
	taskRepo := tasks_repo.NewTaskRepository(dbConn)
	taskService := tasks_usecase.NewTaskService(taskRepo, teamRepo, nil)
	taskHandler := tasks_controller.NewTaskHandler(taskService)

	// 8. Gin роутер
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rate limiter
	redisClient, err := cache.NewCache("redis:6379")
	if err != nil {
		logger.Warn("Redis not available, rate limiter disabled: %v", err)
	} else {
		logger.Info("Redis connected for rate limiter")
		rateLimiter := middleware.NewRateLimiter(redisClient.GetClient(), 100, time.Minute)
		router.Use(middleware.GinRateLimitMiddleware(rateLimiter))
	}

	// Регистрируем роутеры
	auth_controller.RegisterAuthRoutes(router, authHandler)
	teams_controller.RegisterTeamRoutes(router, teamHandler)
	tasks_controller.RegisterTaskRoutes(router, taskHandler, jwtConfig)

	// Метрики
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 9. HTTP сервер
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// 10. Запускаем
	go func() {
		logger.Info("Server started on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error: %v", err)
		}
	}()

	// 11. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error: %v", err)
		return err
	}

	logger.Info("Server stopped")
	return nil
}
