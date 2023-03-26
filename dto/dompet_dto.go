package dto

type DompetCreateDTO struct {
	NamaDompet string `json:"nama_dompet" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
	Saldo      string `json:"saldo" binding:"required"`
}
