package routes

import (
	"dompet-api/controller"
	"dompet-api/middleware"
	"dompet-api/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, userController controller.UserController, dompetController controller.DompetController, jwtService service.JWTService) {
	userRegist := router.Group("/user")
	{
		userRegist.POST("", userController.RegisterUser)
		userRegist.POST("/login", userController.LoginUser)
	}

	userRoutes := router.Group("/secured").Use(middleware.Authenticate())
	{
		userRoutes.GET("/me", dompetController.LihatDompet)
	}
}
