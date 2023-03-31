package dto

type CreatePemasukanDTO struct {
	Deskripsi string `json:"deskripsi" binding:"required"` // berupa detail dari pemasukan/pengeluaran
	Pemasukan uint64 `json:"pemasukan" binding:"required"` // tuliskan nominal pemasukan
	DompetID  uint64 `gorm:"foreignKey" json:"user_id" binding:"required"`
}

type CreatePengeluaranDTO struct {
	Deskripsi   string `json:"deskripsi" binding:"required"`   // berupa detail dari pemasukan/pengeluaran
	Pengeluaran uint64 `json:"pengeluaran" binding:"required"` // tuliskan nominal pemasukan
	DompetID    uint64 `gorm:"foreignKey" json:"user_id" binding:"required"`
}

type TransferRequest struct {
	NamaDompet string `json:"nama_dompet" binding:"required"`
	Nominal    uint64 `json:"nominal" binding:"required"`
	Deskripsi  string `json:"deskripsi" binding:"required"`
	Kategori   string `json:"kategori" binding:"required"`
}
