package controller

import (
	"net/http"
	"dompet-api/dto"
	"dompet-api/entity"
	"dompet-api/service"
	"dompet-api/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type catatanController struct {
	catatanService service.CatatanService
}

type CatatanController interface {
	Transfer(ctx *gin.Context)
	InsertKategori(ctx *gin.Context)
}

func NewCatatanController(cs service.CatatanService) CatatanController {
	return &catatanController{
		catatanService: cs,
	}
}

func (c *catatanController) Transfer(ctx *gin.Context) {

	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	idUser, err := tokenService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	idDompet, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var transferDTO dto.TransferRequest
	errDTO := ctx.ShouldBind(&transferDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	result, err := c.catatanService.Transfer(transferDTO, idUser, idDompet)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berhasil melakukan transfer", http.StatusOK, result)
	ctx.JSON(http.StatusCreated, response)

}

func (c *catatanController) InsertKategori(ctx *gin.Context) {
	var kategori entity.KategoriCatatanKeuangan
	err := ctx.ShouldBind(&kategori)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	result, err := c.catatanService.InsertKategori(kategori)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berhasil insert kategori", http.StatusOK, result)
	ctx.JSON(http.StatusCreated, response)

}
