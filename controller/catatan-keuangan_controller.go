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

	// datepin user id
	userID, err := tokenService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get ID from token", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// dapetin detail pemasukan
	var pemasukanDTO dto.CreatePemasukanDTO
	if tx := ctx.ShouldBind(&pemasukanDTO); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// verif dompet user owner
	verifyDompetUserOwnership, err := c.dompetService.IsDompetOwnedByUserID(ctx.Request.Context(), pemasukanDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verif dompet ke collab user
	verifyDompetUserAccess, err := c.dompetService.IsUserHasAccessToDompet(ctx.Request.Context(), userID, pemasukanDTO.DompetID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify user acces to dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verif
	if verifyDompetUserOwnership != true && verifyDompetUserAccess != true {
		response := utils.BuildErrorResponse("Failed to create pemasukan: user not have access to dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// get dompet
	dompetUpdated, err := c.dompetService.GetDetailDompet(pemasukanDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get dompet for update saldo", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// updet saldo
	dompetUpdated.Saldo += pemasukanDTO.Pemasukan
	_, err = c.dompetService.UpdateDompet(ctx.Request.Context(), dompetUpdated)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to update saldo", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// create pemasukan
	pemasukan, err := c.catatanService.CreatePemasukan(ctx.Request.Context(), pemasukanDTO)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to create pemasukan", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("Success to create pemasukan and update saldo", http.StatusOK, pemasukan)
	ctx.JSON(http.StatusCreated, response)
}

func (c *catatanController) CreatePengeluaran(ctx *gin.Context) {
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

	// get detail pengeluaran
	var pengeluaranDTO dto.CreatePengeluaranDTO
	if tx := ctx.ShouldBind(&pengeluaranDTO); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// verif dompet user owner
	verifyDompetUserOwnership, err := c.dompetService.IsDompetOwnedByUserID(ctx.Request.Context(), pengeluaranDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verif dompet user anggota collab
	verifyDompetUserAccess, err := c.dompetService.IsUserHasAccessToDompet(ctx.Request.Context(), userID, pengeluaranDTO.DompetID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify user acces to dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verif
	if verifyDompetUserOwnership != true && verifyDompetUserAccess != true {
		response := utils.BuildErrorResponse("Failed to create pengeluaran: user not have access to dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// get dompet
	dompetUpdated, err := c.dompetService.GetDetailDompet(pengeluaranDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get dompet for update saldo", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// check saldo
	if dompetUpdated.Saldo < pengeluaranDTO.Pengeluaran {
		response := utils.BuildErrorResponse("Failed to create pengeluaran: saldo tidak cukup", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// update saldo
	dompetUpdated.Saldo -= pengeluaranDTO.Pengeluaran
	_, err = c.dompetService.UpdateDompet(ctx.Request.Context(), dompetUpdated)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to update saldo", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// create catatan pengeluaran
	pengeluaran, err := c.catatanService.CreatePengeluaran(ctx.Request.Context(), pengeluaranDTO)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to create pengeluaran", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("Success to create pengeluaran and update saldo", http.StatusOK, pengeluaran)
	ctx.JSON(http.StatusCreated, response)
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
	verifyDompetUserOwnership, err := c.dompetService.IsDompetOwnedByUserID(ctx.Request.Context(), catatanDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	verifyDompetUserAccess, err := c.dompetService.IsUserHasAccessToDompet(ctx.Request.Context(), userID, catatanDTO.DompetID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify user acces to dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verifyDompetUserOwnership != true && verifyDompetUserAccess != true {
		response := utils.BuildErrorResponse("Failed to delete catatan keuangan: user not have access to dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verify dompet catatan
	verifyDompetCatatan, err := c.catatanService.IsCatatanExistInDompet(ctx.Request.Context(), catatanDTO.ID, catatanDTO.DompetID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet-catatan relationship", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verify user anggota collab dompet
	if verifyDompetCatatan != true {
		response := utils.BuildErrorResponse("Failed to delete catatan keuangan: catatan not exist in dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

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

	// update saldo
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
	verifyDompetUserOwnership, err := c.dompetService.IsDompetOwnedByUserID(ctx.Request.Context(), pemasukanDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	verifyDompetUserAccess, err := c.dompetService.IsUserHasAccessToDompet(ctx.Request.Context(), userID, pemasukanDTO.DompetID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify user acces to dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verify user anggota collab dompet
	if verifyDompetUserOwnership != true && verifyDompetUserAccess != true {
		response := utils.BuildErrorResponse("Failed to update pemasukan: user not have access to dompet", http.StatusBadRequest)
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
		response := utils.BuildErrorResponse("Failed to update catatan pemasukan: catatan not exist in dompet", http.StatusBadRequest)
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

	// update saldo
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

func (c *catatanController) UpdatePengeluaran(ctx *gin.Context) {
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
	var pengeluaranDTO dto.UpdatePengeluaranDTO
	if tx := ctx.ShouldBind(&pengeluaranDTO); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// verivfy dompet user
	verifyDompetUserOwnership, err := c.dompetService.IsDompetOwnedByUserID(ctx.Request.Context(), pengeluaranDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet ownership", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verify user anggota collab dompet
	verifyDompetUserAccess, err := c.dompetService.IsUserHasAccessToDompet(ctx.Request.Context(), userID, pengeluaranDTO.DompetID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify user acces to dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verifyDompetUserOwnership != true && verifyDompetUserAccess != true {
		response := utils.BuildErrorResponse("Failed to update pengeluaran: user not have access to dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// verify dompet catatan
	verifyDompetCatatan, err := c.catatanService.IsCatatanExistInDompet(ctx.Request.Context(), pengeluaranDTO.ID, pengeluaranDTO.DompetID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to verify dompet-catatan relationship", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if verifyDompetCatatan != true {
		response := utils.BuildErrorResponse("Failed to update pengeluaran: catatan not exist in dompet", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	dompetDetail, err := c.dompetService.GetDetailDompet(pengeluaranDTO.DompetID, userID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get dompet for update saldo", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	oldCatatanDetail, err := c.catatanService.GetCatatanByID(ctx.Request.Context(), pengeluaranDTO.ID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to get catatan detail for update saldo", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// update saldo
	dompetDetail.Saldo += oldCatatanDetail.Pengeluaran
	dompetDetail.Saldo -= pengeluaranDTO.Pengeluaran

	_, err = c.catatanService.UpdatePengeluaran(ctx.Request.Context(), pengeluaranDTO)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to update pengeluaran on catatan keuangan", http.StatusBadRequest)
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
