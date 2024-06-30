package main

import (
	"errors"
	"finance-manager-api-service/internal/client/operation_service/category"
	"finance-manager-api-service/internal/client/operation_service/operation"
	"finance-manager-api-service/internal/client/stats_service"
	"finance-manager-api-service/internal/client/user_service"
	"finance-manager-api-service/internal/config"
	"finance-manager-api-service/internal/handler/auth"
	"finance-manager-api-service/internal/handler/operations"
	"finance-manager-api-service/internal/handler/stats"
	"finance-manager-api-service/pkg/cache/freecache"
	"finance-manager-api-service/pkg/jwt"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/metric"
	"finance-manager-api-service/pkg/shutdown"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

func main() {
	logging.InitLogger()
	logger := logging.GetLogger()
	logger.Info("logger initialized")

	logger.Info("config initializing")
	cfg := config.GetConfig()

	logger.Info("router initializing")
	router := httprouter.New()

	logger.Info("cache initializing")
	refreshTokenCache := freecache.NewCacheRepo(104857600) //100MB

	logger.Info("jwt helper initializing")
	jwtHelper := jwt.NewHelper(refreshTokenCache, logger)

	logger.Info("create and register handlers")

	metricHandler := metric.NewHandler(logger)
	metricHandler.Register(router)

	userService := user_service.NewService(cfg.UserService.URL, "/users", logger)
	authHandler := auth.NewAuthHandler(logger, userService, jwtHelper)
	authHandler.Register(router)

	categoryService := category.NewService(cfg.OperationService.URL, "/categories", logger)
	categoryHandler := operations.NewCategoryHandler(logger, categoryService)
	categoryHandler.Register(router)

	operationService := operation.NewService(cfg.OperationService.URL, "/operations", logger)
	operationHandler := operations.NewOperationHandler(logger, operationService)
	operationHandler.Register(router)

	statsService := stats_service.NewService(cfg.StatsService.URL, "/stats", logger)
	statsHandler := stats.NewHandler(logger, statsService)
	statsHandler.Register(router)

	logger.Info("start application")
	start(router, logger, cfg)
}

func start(router *httprouter.Router, logger *logging.Logger, cfg *config.Config) {
	logger.Infof("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	logger.Info("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
