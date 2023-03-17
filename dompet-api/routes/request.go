package routes

import (
	"oprec/dompet-api/controller"
	"oprec/dompet-api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userController controller.UserController, dompetController controller.DompetController) {
	userRegist := router.Group("/masuk")
	{
		userRegist.POST("/register", userController.Register)
		userRegist.POST("/login", userController.Login)
	}

	userRoutes := router.Group("/secured").Use(middleware.Authenticate())
	{
		userRoutes.GET("/me", dompetController.LihatDompet)
	}

}
