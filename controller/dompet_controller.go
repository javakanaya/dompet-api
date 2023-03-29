package controller

import (
	"dompet-api/dto"
	"dompet-api/service"
	"dompet-api/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type dompetController struct {
	dompetService service.DompetService
}

type DompetController interface {
	LihatDompet(ctx *gin.Context)
	CreateDompet(ctx *gin.Context)
	DetailDompet(ctx *gin.Context)
	Invite(ctx *gin.Context)
	DeleteDompet(ctx *gin.Context)
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

func (c *dompetController) DeleteDompet(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	userID, err := tokenService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get ID from token", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	dompetID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("Gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	verify, err := c.dompetService.IsDompetOwnedByUserID(ctx, dompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verify == true {
		err = c.dompetService.DeleteDompet(ctx.Request.Context(), dompetID)
		if err != nil {
			response := utils.BuildErrorResponse("Failed to delete dompet", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		response := utils.BuildResponse("Success to delete dompet", http.StatusOK, nil)
		ctx.JSON(http.StatusCreated, response)
		return
	}

	response := utils.BuildErrorResponse("Failed to delete: wrong dompet ownership", http.StatusBadRequest)
	ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
}
