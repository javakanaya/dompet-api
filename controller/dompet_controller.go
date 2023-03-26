package controller

import (
	"dompet-api/dto"
	"dompet-api/service"
	"dompet-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type dompetController struct {
	dompetService service.DompetService
}

type DompetController interface {
	LihatDompet(ctx *gin.Context)
	CreateDompet(ctx *gin.Context)
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

func (c *dompetController) CreateDompet(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	id, err := tokenService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get ID from token", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var dompet dto.DompetCreateDTO
	if tx := ctx.ShouldBind(&dompet); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	dompet.UserID = id

	result, err := c.dompetService.CreateDompet(ctx.Request.Context(), dompet)
	if err != nil {
		res := utils.BuildErrorResponse("Failed to create dompet", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse("Success to  create dompet", http.StatusOK, result)
	ctx.JSON(http.StatusOK, res)
}
