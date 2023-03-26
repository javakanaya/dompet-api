package dto

type DompetCreateDTO struct {
	NamaDompet string `json:"nama_dompet" binding:"required"`
	UserID     uint64 `json:"user_id"`
	Saldo      uint64 `json:"saldo" binding:"required"`
}