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
	}
}
