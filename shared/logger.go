package shared

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger 创建新的日志记录器
func NewLogger() *zap.Logger {
	config := GetLoggingConfig()
	
	// 设置日志级别
	var level zapcore.Level
	switch config.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	// 创建编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 设置输出
	var cores []zapcore.Core
	
	// 控制台输出
	consoleCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
	cores = append(cores, consoleCore)

	// 文件输出（如果配置了）
	if config.File != "" {
		file, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			fileCore := zapcore.NewCore(
				encoder,
				zapcore.AddSync(file),
				level,
			)
			cores = append(cores, fileCore)
		}
	}

	// 合并所有cores
	core := zapcore.NewTee(cores...)

	// 创建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}

// NewProductionLogger 创建生产环境日志记录器（简化版本）
func NewProductionLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

// NewDevelopmentLogger 创建开发环境日志记录器（简化版本）
func NewDevelopmentLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

