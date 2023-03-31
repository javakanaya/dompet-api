package dto

type DompetCreateDTO struct {
	NamaDompet string  `json:"nama_dompet" binding:"required"`
	UserID     uint64  `json:"user_id"`
	Saldo      *uint64 `json:"saldo" binding:"required"`
}

type InviteUserRequest struct {
	DompetID  uint64 `json:"dompet_id"`
	EmailUser string `json:"user_email" binding:"required"`
}

type DompetUpdateSaldoDTO struct {
	NamaDompet string  `json:"nama_dompet" binding:"required"`
	Saldo      *uint64 `json:"saldo" binding:"required"`
}
