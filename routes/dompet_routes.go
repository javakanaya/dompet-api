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
		dompetRoutes.POST("/dompet", dompetController.CreateDompet)
		dompetRoutes.GET("/dompet/:id", dompetController.DetailDompet)
		dompetRoutes.DELETE("/dompet/:id", dompetController.DeleteDompet)
		dompetRoutes.PUT("/collab/:id", dompetController.Invite)
		dompetRoutes.PUT("/update/:id", dompetController.UpdateNama)
	}
}
