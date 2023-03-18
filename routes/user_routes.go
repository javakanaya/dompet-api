package routes

import (
	"dompet-api/controller"
	"dompet-api/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("", userController.RegisterUser)
		userRoutes.POST("/login", userController.LoginUser)
	}
}
