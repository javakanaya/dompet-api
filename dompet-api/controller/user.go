package controller

import (
	"net/http"
	"oprec/dompet-api/dto"
	"oprec/dompet-api/service"
	"oprec/dompet-api/utils"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
}

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

func NewUserController(us service.UserService) UserController {
	return &userController{
		userService: us,
	}
}

func (c *userController) Register(ctx *gin.Context) {
	var userDTO dto.UserRegisterRequest

	errDTO := ctx.ShouldBind(&userDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := c.userService.CreateUser(ctx, userDTO)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponse("Akun catatan keuangan anda berhasil dibuat", http.StatusCreated, user)
	ctx.JSON(http.StatusCreated, response)

}

func (c *userController) Login(ctx *gin.Context) {
	var userDTO dto.UserLoginRequest

	errDTO := ctx.ShouldBind(&userDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := c.userService.FindUserByEmail(ctx, userDTO.Email)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	checkPass, _ := utils.ComparePassword(user.Password, []byte(userDTO.Password))
	if !checkPass {
		response := utils.BuildErrorResponse("Password Salah", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	tokenService := service.NewJWTService()
	tokenString := tokenService.GenerateToken(user.ID, user.Name)

	response := utils.BuildResponse("Anda berhasil login, berikut token anda:", http.StatusCreated, tokenString)
	ctx.JSON(http.StatusCreated, response)

}
