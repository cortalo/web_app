package handler

import (
	"net/http"
	"web_app/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService domain.UserService // 依赖接口，不依赖具体类型
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var p ParamRegister
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Info("register with invalid param", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid JSON format", "error": err.Error()})
		return
	}
	if err := h.userService.Register(p.Username, p.Password); err != nil {
		c.JSON(httpStatusFromError(err), gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h *UserHandler) Login(c *gin.Context) {
	var p ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Info("login with invalid param", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid JSON format", "error": err.Error()})
		return
	}
	if err := h.userService.Login(p.Username, p.Password); err != nil {
		zap.L().Info("login with error", zap.Error(err))
		c.JSON(httpStatusFromError(err), gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "success")
}
