package entity

type CatatanKeuangan struct {
	ID          uint64  `json:"id" gorm:"primaryKey"`
	Deskripsi   string  `json:"deskripsi" binding:"required"`   // berupa detail dari pemasukan/pengeluaran
	Pemasukan   uint64  `json:"pemasukan" binding:"required"`   // tuliskan nominal pemasukan
	Pengeluaran uint64  `json:"pengeluaran" binding:"required"` // tuliskan nominal pengeluaran
	DompetID    uint64  `gorm:"foreignKey" json:"user_id"`
	Dompet      *Dompet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}

func (CatatanKeuangan) TableName() string {
	return "catatan_keuangan"
}
