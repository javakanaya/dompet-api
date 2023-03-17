package controller

import (
	"net/http"
	"oprec/dompet-api/service"
	"oprec/dompet-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type dompetController struct {
	dompetService service.DompetService
}

type DompetController interface {
	LihatDompet(ctx *gin.Context)
}

func NewDompetController(dc service.DompetService) DompetController {
	return &dompetController{
		dompetService: dc,
	}
}

func (c *dompetController) LihatDompet(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	id, _ := tokenService.GetUserIDByToken(token)

	getDompet, err := c.dompetService.GetMyDompet(id)
	if err != nil {
		response := utils.BuildErrorResponse("Gagal melihat dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berhasil", http.StatusOK, getDompet)
	ctx.JSON(http.StatusOK, response)
}
