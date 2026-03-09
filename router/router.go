package router

import (
	"web_app/handler"
	"web_app/logger"

	"github.com/gin-gonic/gin"
)

func Setup(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger())
	r.Use(logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.Register)
		}
		login := v1.Group("/login")
		{
			login.POST("", userHandler.Login)
		}
	}

	return r
}
