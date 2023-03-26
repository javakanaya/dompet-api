package routes

import (
	"dompet-api/controller"
	"dompet-api/middleware"
	"dompet-api/service"

	"github.com/gin-gonic/gin"
)

func DompetRouter(router *gin.Engine, dompetController controller.DompetController, jwtService service.JWTService) {
	dompetRoutes := router.Group("/secured").Use(middleware.Authenticate())
	{
		dompetRoutes.GET("/me", dompetController.LihatDompet)
		dompetRoutes.POST("/create/dompet", dompetController.BuatDompet)
		dompetRoutes.GET("/dompet/:id", dompetController.DetailDompet)
		dompetRoutes.PUT("/collab/:id", dompetController.Invite)
	}
}
