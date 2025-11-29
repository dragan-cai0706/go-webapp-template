package shared

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// RespondError 返回错误响应
func RespondError(c *gin.Context, statusCode int, message string, err error, logger *zap.Logger) {
	if logger != nil {
		if err != nil {
			logger.Error(message, zap.Error(err), zap.Int("status", statusCode))
		} else {
			logger.Warn(message, zap.Int("status", statusCode))
		}
	}

	response := ErrorResponse{
		Code:    statusCode,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(statusCode, response)
}

// RespondSuccess 返回成功响应
func RespondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	})
}

// RespondBadRequest 返回400错误
func RespondBadRequest(c *gin.Context, message string, err error, logger *zap.Logger) {
	RespondError(c, http.StatusBadRequest, message, err, logger)
}

// RespondInternalError 返回500错误
func RespondInternalError(c *gin.Context, message string, err error, logger *zap.Logger) {
	RespondError(c, http.StatusInternalServerError, message, err, logger)
}

// FormatServerAddr 格式化服务器地址
func FormatServerAddr(port int) string {
	return fmt.Sprintf(":%d", port)
}

