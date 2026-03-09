package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 统一处理参数绑定
func bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		zap.L().Info("invalid param", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid param", "error": err.Error()})
		return false
	}
	return true
}

// 统一处理响应
func responseError(c *gin.Context, err error) {
	zap.L().Info("request error", zap.Error(err))
	c.JSON(httpStatusFromError(err), gin.H{"msg": err.Error()})
}

func responseSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"msg": "success", "data": data})
}
