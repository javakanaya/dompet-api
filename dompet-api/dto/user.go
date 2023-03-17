package dto

type UserRegisterRequest struct { // seseorang ingin melakukan registrasi akun catatan keuangan
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct { // seseorang ingin melakukan login ke akun
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
