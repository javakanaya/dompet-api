package controller

import (
	"dompet-api/dto"
	"dompet-api/service"
	"dompet-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type catatanController struct {
	catatanService service.CatatanService
	dompetService  service.DompetService
}

type CatatanController interface {
	CreatePemasukan(ctx *gin.Context)
	// CreatePengeluaran(ctx *gin.Context)
}

func NewCatatanController(cs service.CatatanService, ds service.DompetService) CatatanController {
	return &catatanController{
		catatanService: cs,
		dompetService:  ds,
	}
}

func (c *catatanController) CreatePemasukan(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	userID, err := tokenService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get ID from token", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var pemasukan dto.CreatePemasukanDTO
	if tx := ctx.ShouldBind(&pemasukan); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	verify, err := c.dompetService.IsDompetOwnedByUserID(ctx, pemasukan.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verify == true {
		pemasukan, err := c.catatanService.CreatePemasukan(ctx.Request.Context(), pemasukan)
		if err != nil {
			response := utils.BuildErrorResponse("Failed to create pemasukan", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		dompetUpdated, err := c.dompetService.GetDetailDompet(pemasukan.DompetID, userID)
		dompetUpdated.Saldo += pemasukan.Pemasukan

		dompetUpdated, err = c.dompetService.UpdateDompet(ctx.Request.Context(), dompetUpdated)
		if err != nil {
			response := utils.BuildErrorResponse("Failed to update saldo", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		response := utils.BuildResponse("Success to create pemasukan and update saldo", http.StatusOK, pemasukan)
		ctx.JSON(http.StatusCreated, response)
		return
	}

	response := utils.BuildErrorResponse("Failed to create pemasukan: wrong dompet ownership", http.StatusBadRequest)
	ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
}
