package handler

import (
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
	var p ParamRegister
	if !bindJSON(c, &p) {
		return
	}
	if err := h.userService.Register(p.Username, p.Password); err != nil {
		responseError(c, err)
		return
	}
	responseSuccess(c, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var p ParamLogin
	if !bindJSON(c, &p) {
		return
	}
	if err := h.userService.Login(p.Username, p.Password); err != nil {
		responseError(c, err)
		return
	}
	responseSuccess(c, nil)
}
