package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"freecharge/go-freecharge/handler"
	"freecharge/go-freecharge/middleware"
	"freecharge/go-freecharge/shared"
)

func main() {
	// 加载配置（shared包会自动加载.env）
	config := shared.GetConfig()

	// 初始化日志
	logger := shared.NewLogger()
	defer logger.Sync()

	logger.Info("服务启动中...")

	// 初始化Redis客户端（必需）
	redisClient := shared.NewRedisClient(config.Redis)
	logger.Info("Redis客户端已初始化", 
		zap.String("addr", config.Redis.Addr),
		zap.String("keyPrefix", config.Redis.KeyPrefix),
	)

	// 初始化处理器
	healthHandler := handler.NewHealthHandler(logger)

	// 设置路由
	router := setupRouter(healthHandler, logger, config)

	// 创建HTTP服务器
	serverAddr := shared.FormatServerAddr(config.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		logger.Info("服务启动", zap.String("addr", serverAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	// 等待关闭信号
	waitForShutdown(server, logger, redisClient)
}

// setupRouter 设置路由
func setupRouter(
	healthHandler *handler.HealthHandler,
	logger *zap.Logger,
	config *shared.Config,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 添加中间件
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.ErrorHandlingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())

	// 健康检查路由
	api := router.Group("/api/")
	{
		api.GET("/health", healthHandler.Health)
		api.GET("/health/readiness", healthHandler.Readiness)
		api.GET("/health/liveness", healthHandler.Liveness)
	}

	return router
}

// waitForShutdown 等待关闭信号并优雅关闭
func waitForShutdown(server *http.Server, logger *zap.Logger, redisClient *shared.RedisClient) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("收到关闭信号，正在优雅关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("服务关闭失败", zap.Error(err))
	} else {
		logger.Info("服务已优雅关闭")
	}

	// 关闭Redis连接
	if err := redisClient.Close(); err != nil {
		logger.Error("Redis连接关闭失败", zap.Error(err))
	} else {
		logger.Info("Redis连接已关闭")
	}
}

