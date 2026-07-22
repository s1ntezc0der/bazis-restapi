package run

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	auth_controller "github.com/s1ntezc0der/bazis-restapi/internal/services/auth/controller"
	auth_repo "github.com/s1ntezc0der/bazis-restapi/internal/services/auth/repository"
	auth_usecase "github.com/s1ntezc0der/bazis-restapi/internal/services/auth/usecase"
	teams_controller "github.com/s1ntezc0der/bazis-restapi/internal/services/teams/controller"
	teams_repo "github.com/s1ntezc0der/bazis-restapi/internal/services/teams/repository"
	teams_usecase "github.com/s1ntezc0der/bazis-restapi/internal/services/teams/usecase"
	tasks_controller "github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/controller"
	tasks_repo "github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/repository"
	tasks_usecase "github.com/s1ntezc0der/bazis-restapi/internal/services/tasks/usecase"
	"github.com/s1ntezc0der/bazis-restapi/internal/config"
	"github.com/s1ntezc0der/bazis-restapi/pkg/cache"
	"github.com/s1ntezc0der/bazis-restapi/pkg/db"
	"github.com/s1ntezc0der/bazis-restapi/pkg/jwt"
	"github.com/s1ntezc0der/bazis-restapi/pkg/logger"
	"github.com/s1ntezc0der/bazis-restapi/pkg/metrics"
	"github.com/s1ntezc0der/bazis-restapi/pkg/middleware"
)

func Run() error {
	// 1. Загружаем конфиг
	cfg := config.Load()

	// 2. Инициализируем логгер
	logger := logger.New()

	// 3. Подключаемся к MySQL
	dbConn, err := db.NewMySQLDB(cfg.DB)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	// 4. Инициализируем JWT
	jwtConfig := jwt.NewJWT(cfg.JWT.Secret, time.Duration(cfg.JWT.Expire)*time.Hour)

	// 5. Auth — Repository → Usecase → Handler
	authRepo := auth_repo.NewAuthRepository(dbConn)
	authService := auth_usecase.NewAuthService(authRepo, jwtConfig)
	authHandler := auth_controller.NewAuthHandler(authService)

	// 6. Teams — Repository → Usecase → Handler
	teamRepo := teams_repo.NewTeamRepository(dbConn)
	teamService := teams_usecase.NewTeamService(teamRepo, authRepo)
	teamHandler := teams_controller.NewTeamHandler(teamService)

	// 7. Tasks — Repository → Usecase → Handler
	taskRepo := tasks_repo.NewTaskRepository(dbConn)
	taskService := tasks_usecase.NewTaskService(taskRepo, authRepo, teamRepo)
	taskHandler := tasks_controller.NewTaskHandler(taskService)

	// 8. Собираем роутеры
	authRouter := auth_controller.NewAuthRouter(authHandler)
	teamRouter := teams_controller.NewTeamRouter(teamHandler)
	taskRouter := tasks_controller.NewTaskRouter(taskHandler)

	// 9. Главный роутер с middleware
	mainRouter := chi.NewRouter()

	// Middleware
	mainRouter.Handle("/metrics", promhttp.Handler())
	mainRouter.Use(middleware.LoggingMiddleware)
	mainRouter.Use(middleware.Recoverer)
	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Монтируем все роутеры
	mainRouter.Mount("/", authRouter)   // /api/v1/register, /api/v1/login
	mainRouter.Mount("/", teamRouter)   // /api/v1/teams, /api/v1/teams/{id}/invite
	mainRouter.Mount("/", taskRouter)   // /api/v1/tasks, /api/v1/tasks/{id}/history

	redisClient, _ := cache.NewCache(cfg.RedisAddr)

	rateLimiter := middleware.NewRateLimiter(redisClient, 100, time.Minute)
	mainRouter.Use(middleware.RateLimitMiddleware(rateLimiter))

	// 10. HTTP сервер
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mainRouter,
	}

	// 11. Запускаем сервер в горутине
	go func() {
		logger.Info("Server started on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error: %v", err)
		}
	}()

	// 12. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit, 
		syscall.SIGINT, 
		syscall.SIGTERM,
	)
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

