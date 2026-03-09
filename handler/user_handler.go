package handler

import (
	"net/http"
	"web_app/domain"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domain.UserService // 依赖接口，不依赖具体类型
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	c.JSON(http.StatusOK, "hello") // 先空着
}
