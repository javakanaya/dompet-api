package routes

import (
	"dompet-api/controller"
	"dompet-api/middleware"
	"dompet-api/service"

	"github.com/gin-gonic/gin"
)

func CatatanRouter(router *gin.Engine, catatanController controller.CatatanController, jwtService service.JWTService) {
	catatanRoutes := router.Group("/secured").Use(middleware.Authenticate())
	{
		catatanRoutes.POST("/transfer/:id", catatanController.Transfer)
		catatanRoutes.POST("/pemasukan", catatanController.CreatePemasukan)
		catatanRoutes.POST("/pengeluaran", catatanController.CreatePengeluaran)
		catatanRoutes.DELETE("/catatan", catatanController.DeleteCatatan)
	}
	
	publicRoutes := router.Group("")
	{
		publicRoutes.GET("/kategori/:jenis", catatanController.Kategori)
	}
	
	adminRoutes := router.Group("/admin")
	{
		adminRoutes.POST("/kategori", catatanController.InsertKategori)
	}
}
