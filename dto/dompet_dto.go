package dto

type DompetCreateRequest struct {
	NamaDompet string `json:"nama_dompet" binding:"required"`
	Saldo      uint64 `json:"saldo" binding:"required"`
	UserID     uint64 `json:"user_id"`
}

type InviteUserRequest struct {
	DompetID  uint64 `json:"dompet_id"`
	EmailUser string `json:"user_email" binding:"required"`
}
