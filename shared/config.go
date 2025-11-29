package shared

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// LoadEnv 加载环境变量
func LoadEnv() {
	var envFile string
	var loaded bool

	// 首先尝试从可执行文件同目录加载
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("警告: 无法获取可执行文件路径: %v", err)
	} else {
		// 获取可执行文件所在目录
		execDir := filepath.Dir(execPath)
		envFile = filepath.Join(execDir, ".env")

		// 检查.env文件是否存在
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				log.Printf("警告: 无法加载.env文件 %s: %v", envFile, err)
			} else {
				log.Printf("成功从可执行文件目录加载环境变量文件: %s", envFile)
				loaded = true
			}
		}
	}

	// 如果可执行文件目录没有找到.env文件，尝试当前工作目录
	if !loaded {
		workDir, err := os.Getwd()
		if err != nil {
			log.Printf("警告: 无法获取当前工作目录: %v", err)
		} else {
			envFile = filepath.Join(workDir, ".env")

			// 检查.env文件是否存在
			if _, err := os.Stat(envFile); err == nil {
				if err := godotenv.Load(envFile); err != nil {
					log.Printf("警告: 无法加载.env文件 %s: %v", envFile, err)
				} else {
					log.Printf("成功从当前工作目录加载环境变量文件: %s", envFile)
					loaded = true
				}
			}
		}
	}

	// 如果两个位置都没有找到.env文件，尝试项目根目录
	if !loaded {
		_ = godotenv.Load(".env")
		_ = godotenv.Overload("../.env")
	}

	// 如果所有位置都没有找到.env文件
	if !loaded {
		log.Println("未找到.env文件（已检查可执行文件目录和当前工作目录），使用系统环境变量")
	}
}

// init 包导入时自动加载.env文件
func init() {
	LoadEnv()
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           // 服务器端口
	ReadTimeout  time.Duration // 读取超时
	WriteTimeout time.Duration // 写入超时
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string // Redis服务器地址
	Password string // Redis密码
	DB       int    // Redis数据库编号
	KeyPrefix string // Redis键前缀
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level   string // 日志级别 (debug, info, warn, error)
	File    string // 日志文件路径
	MaxSize int    // 日志文件最大大小(MB)
}

// Config 应用程序配置
type Config struct {
	Server  ServerConfig
	Redis   RedisConfig
	Logging LoggingConfig
}

// GetEnv 获取环境变量字符串
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt 获取环境变量整数
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetEnvDuration 获取环境变量时长
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetEnvBool 获取环境变量布尔值
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:         GetEnvInt("GO_SERVER_PORT", 8080),
		ReadTimeout:  GetEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: GetEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
	}
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     GetEnv("REDIS_ADDR", "127.0.0.1:6379"),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       GetEnvInt("REDIS_DB", 0),
		KeyPrefix: GetEnv("REDIS_KEY_PREFIX", ""),
	}
}

// GetLoggingConfig 获取日志配置
func GetLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		Level:   GetEnv("LOG_LEVEL", "info"),
		File:    GetEnv("LOG_FILE", ""),
		MaxSize: GetEnvInt("LOG_MAX_SIZE", 100),
	}
}

// GetConfig 获取完整应用程序配置
func GetConfig() *Config {
	return &Config{
		Server:  *GetServerConfig(),
		Redis:   *GetRedisConfig(),
		Logging: *GetLoggingConfig(),
	}
}

