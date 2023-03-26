package controller

import (
	"net/http"
	"dompet-api/dto"
	"dompet-api/service"
	"dompet-api/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type dompetController struct {
	dompetService service.DompetService
}

type DompetController interface {
	LihatDompet(ctx *gin.Context)
	BuatDompet(ctx *gin.Context)
	DetailDompet(ctx *gin.Context)
	Invite(ctx *gin.Context)
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

func (c *dompetController) BuatDompet(ctx *gin.Context) { // ini temporary, buat dompetnya masih pake punya aku jav
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1) // menghilangkan kata bearer dari token, karena mau di proses
	tokenService := service.NewJWTService()

	id, _ := tokenService.GetUserIDByToken(token) // id sudah pasti ada jika berhasil melewati validate, jadi tidak mungkin error

	var dompetDTO dto.DompetCreateRequest
	errDTO := ctx.ShouldBind(&dompetDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newDTO := dto.DompetCreateRequest{
		NamaDompet: dompetDTO.NamaDompet,
		Saldo:      dompetDTO.Saldo,
		UserID:     id,
	}

	activeDompet, err := c.dompetService.CreateDompet(newDTO)
	if err != nil {
		response := utils.BuildErrorResponse("gagal membuat dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("dompet anda berhasil dibuat", http.StatusCreated, activeDompet)
	ctx.JSON(http.StatusCreated, response)

}

func (c *dompetController) DetailDompet(ctx *gin.Context) {
	idDompet, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	result, err := c.dompetService.GetDetailDompet(idDompet)
	if err != nil {
		response := utils.BuildErrorResponse("invalid ID", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berikut detail dari dompet:", http.StatusOK, result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *dompetController) Invite(ctx *gin.Context) {
	idDompet, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var inviteDTO dto.InviteUserRequest
	errDTO := ctx.ShouldBind(&inviteDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newDTO := dto.InviteUserRequest{
		DompetID:  idDompet,
		EmailUser: inviteDTO.EmailUser,
	}

	result, err := c.dompetService.InviteToDompet(newDTO)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("berhasil melakukan penambahan partisipan", http.StatusOK, result)
	ctx.JSON(http.StatusCreated, response)
}
