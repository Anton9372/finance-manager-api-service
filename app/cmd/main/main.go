package main

import (
	"errors"
	_ "finance-manager-api-service/docs"
	"finance-manager-api-service/internal/client/operation_service/category"
	"finance-manager-api-service/internal/client/operation_service/operation"
	"finance-manager-api-service/internal/client/stats_service"
	"finance-manager-api-service/internal/client/user_service"
	user_service_grpc "finance-manager-api-service/internal/client/user_service/grpc/v1"
	"finance-manager-api-service/internal/client/user_service/http"
	"finance-manager-api-service/internal/config"
	"finance-manager-api-service/internal/handler/auth"
	"finance-manager-api-service/internal/handler/categories"
	"finance-manager-api-service/internal/handler/operations"
	"finance-manager-api-service/internal/handler/stats"
	"finance-manager-api-service/internal/handler/users"
	"finance-manager-api-service/pkg/cache/freecache"
	"finance-manager-api-service/pkg/jwt"
	"finance-manager-api-service/pkg/logging"
	"finance-manager-api-service/pkg/metric"
	"finance-manager-api-service/pkg/shutdown"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

// @Title		Finance-manager API
// @Version		1.0
// @Description	Finance-manager application

// @Contact.name	Anton
// @Contact.email	ap363402@gmail.com

// @License.name Apache 2.0

// @Host 		localhost:10000
// @BasePath 	/api
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

	logger.Info("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	metricHandler := metric.NewHandler(logger)
	metricHandler.Register(router)

	var userService user_service.UserService
	if cfg.UserService.ConnectWithGRPC == true {
		var err error
		logger.Info("connect to user service through grpc")
		userService, err = user_service_grpc.NewClient(cfg.UserService.GrpcUrl, logger)
		if err != nil {
			logger.Fatal(err.Error())
		}
	} else {
		logger.Info("connect to user service through http")
		userService = user_service_http.NewService(cfg.UserService.HttpUrl, "/users", logger)
	}
	authHandler := auth.NewAuthHandler(logger, userService, jwtHelper)
	authHandler.Register(router)
	userHandler := users.NewUserHandler(logger, userService)
	userHandler.Register(router)

	categoryService := category.NewService(cfg.OperationService.URL, "/categories", logger)
	categoryHandler := categories.NewCategoryHandler(logger, categoryService)
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
	logger.Infof("bind application to host: %s and port: %d", cfg.HTTP.IP, cfg.HTTP.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.HTTP.IP, cfg.HTTP.Port))
	if err != nil {
		logger.Fatal(err)
	}

	c := cors.New(cors.Options{
		AllowedMethods:   cfg.HTTP.CORS.AllowedMethods,
		AllowedOrigins:   cfg.HTTP.CORS.AllowedOrigins,
		AllowCredentials: cfg.HTTP.CORS.AllowCredentials,
		AllowedHeaders:   cfg.HTTP.CORS.AllowedHeaders,
		ExposedHeaders:   cfg.HTTP.CORS.ExposedHeaders,
	})

	handler := c.Handler(router)

	server := &http.Server{
		Handler: handler,
		//TODO to config
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
