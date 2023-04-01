package controller

import (
	"dompet-api/dto"
	"dompet-api/entity"
	"dompet-api/service"
	"dompet-api/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type catatanController struct {
	catatanService service.CatatanService
	dompetService  service.DompetService
}

type CatatanController interface {
	Transfer(ctx *gin.Context)
	InsertKategori(ctx *gin.Context)
	CreatePemasukan(ctx *gin.Context)
	CreatePengeluaran(ctx *gin.Context)
	DeleteCatatan(ctx *gin.Context)
	UpdatePemasukan(ctx *gin.Context)
	UpdatePengeluaran(ctx *gin.Context)
}

func NewCatatanController(cs service.CatatanService, ds service.DompetService) CatatanController {
	return &catatanController{
		catatanService: cs,
		dompetService:  ds,
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

	var pemasukanDTO dto.CreatePemasukanDTO
	if tx := ctx.ShouldBind(&pemasukanDTO); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	verify, err := c.dompetService.IsDompetOwnedByUserID(ctx, pemasukanDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verify == true {
		dompetUpdated, err := c.dompetService.GetDetailDompet(pemasukanDTO.DompetID, userID)
		if err != nil {
			response := utils.BuildErrorResponse("Failed to get dompet for update saldo", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		pemasukan, err := c.catatanService.CreatePemasukan(ctx.Request.Context(), pemasukanDTO)
		if err != nil {
			response := utils.BuildErrorResponse("Failed to create pemasukan", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		dompetUpdated.Saldo += pemasukan.Pemasukan

		_, err = c.dompetService.UpdateDompet(ctx.Request.Context(), dompetUpdated)
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

func (c *catatanController) CreatePengeluaran(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	userID, err := tokenService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get ID from token", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	var pengeluaranDTO dto.CreatePengeluaranDTO
	if tx := ctx.ShouldBind(&pengeluaranDTO); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	verify, err := c.dompetService.IsDompetOwnedByUserID(ctx, pengeluaranDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verify == true {
		dompetUpdated, err := c.dompetService.GetDetailDompet(pengeluaranDTO.DompetID, userID)
		if dompetUpdated.Saldo < pengeluaranDTO.Pengeluaran {
			response := utils.BuildErrorResponse("Failed to create pengeluaran: saldo tidak cukup", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		dompetUpdated.Saldo -= pengeluaranDTO.Pengeluaran

		pengeluaran, err := c.catatanService.CreatePengeluaran(ctx.Request.Context(), pengeluaranDTO)
		if err != nil {
			response := utils.BuildErrorResponse("Failed to create pengeluaran", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		_, err = c.dompetService.UpdateDompet(ctx.Request.Context(), dompetUpdated)
		if err != nil {
			response := utils.BuildErrorResponse("Failed to update saldo", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		response := utils.BuildResponse("Success to create pengeluaran and update saldo", http.StatusOK, pengeluaran)
		ctx.JSON(http.StatusCreated, response)
		return
	}

	response := utils.BuildErrorResponse("Failed to create pengeluaran: wrong dompet ownership", http.StatusBadRequest)
	ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
}

func (c *catatanController) DeleteCatatan(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	// get user ID
	userID, err := tokenService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get ID from token", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// dapetin ID catatan dan Dompet ID yang mau dihapus
	var catatanDTO dto.DeleteCatatanDTO
	if tx := ctx.ShouldBind(&catatanDTO); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// verivfy dompet user
	verifyDompetUser, err := c.dompetService.IsDompetOwnedByUserID(ctx.Request.Context(), catatanDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verifyDompetUser == true {
		// verify dompet catatan
		verifyDompetCatatan, err := c.catatanService.IsCatatanExistInDompet(ctx.Request.Context(), catatanDTO.ID, catatanDTO.DompetID)
		if err != nil {
			response := utils.BuildErrorResponse("Failed to verify dompet-catatan relationship", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		if verifyDompetCatatan == true {

			dompetUpdated, err := c.dompetService.GetDetailDompet(catatanDTO.DompetID, userID)
			if err != nil {
				response := utils.BuildErrorResponse("Failed to get dompet for update saldo", http.StatusBadRequest)
				ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
				return
			}

			catatanDetail, err := c.catatanService.GetCatatanByID(ctx.Request.Context(), catatanDTO.ID)
			if err != nil {
				response := utils.BuildErrorResponse("Failed to get catatan detail for update saldo", http.StatusBadRequest)
				ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
				return
			}

			dompetUpdated.Saldo += catatanDetail.Pengeluaran
			dompetUpdated.Saldo -= catatanDetail.Pemasukan

			_, err = c.dompetService.UpdateDompet(ctx.Request.Context(), dompetUpdated)
			if err != nil {
				response := utils.BuildErrorResponse("Failed to update saldo", http.StatusBadRequest)
				ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
				return
			}

			err = c.catatanService.DeleteCatatanKeuangan(ctx.Request.Context(), catatanDTO.ID)
			if err != nil {
				response := utils.BuildErrorResponse("Failed to delete catatan keuangan", http.StatusBadRequest)
				ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
				return
			}
			response := utils.BuildResponse("Success to delete catatan keuangan", http.StatusOK, nil)
			ctx.JSON(http.StatusCreated, response)
			return
		}
		response := utils.BuildErrorResponse("Failed to delete catatan keuangan: catatan not exist in dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	response := utils.BuildErrorResponse("Failed to delete catatan keuangan: wrong dompet ownership", http.StatusBadRequest)
	ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
}

func (c *catatanController) UpdatePemasukan(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	// get user ID
	userID, err := tokenService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get ID from token", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// dapetin ID catatan dan Dompet ID yang mau dihapus
	var pemasukanDTO dto.UpdatePemasukanDTO
	if tx := ctx.ShouldBind(&pemasukanDTO); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// verivfy dompet user
	verifyDompetUser, err := c.dompetService.IsDompetOwnedByUserID(ctx.Request.Context(), pemasukanDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verifyDompetUser != true {
		response := utils.BuildErrorResponse("Failed to delete catatan keuangan: wrong dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verify dompet catatan
	verifyDompetCatatan, err := c.catatanService.IsCatatanExistInDompet(ctx.Request.Context(), pemasukanDTO.ID, pemasukanDTO.DompetID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet-catatan relationship", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verifyDompetCatatan != true {
		response := utils.BuildErrorResponse("Failed to delete catatan keuangan: catatan not exist in dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	dompetDetail, err := c.dompetService.GetDetailDompet(pemasukanDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get dompet for update saldo", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	oldCatatanDetail, err := c.catatanService.GetCatatanByID(ctx.Request.Context(), pemasukanDTO.ID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get catatan detail for update saldo", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	dompetDetail.Saldo -= oldCatatanDetail.Pemasukan
	dompetDetail.Saldo += pemasukanDTO.Pemasukan

	_, err = c.catatanService.UpdatePemasukan(ctx.Request.Context(), pemasukanDTO)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to update pemasukan on catatan keuangan", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	_, err = c.dompetService.UpdateDompet(ctx.Request.Context(), dompetDetail)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to update saldo on dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
}
