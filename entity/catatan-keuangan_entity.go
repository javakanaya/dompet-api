package entity

type CatatanKeuangan struct {
	ID           uint64  `json:"id" gorm:"primaryKey"`
	JenisCatatan string  `json:"jenis_catatan" binding:"required"` // berupa pengeluaran, pemasukan, transfer atau apa
	Nominal      uint64  `json:"nominal" binding:"required"`       // tuliskan nominal pengeluaran, pemasukannya
	DompetID     uint64  `gorm:"foreignKey" json:"user_id"`
	Dompet       *Dompet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}

func (CatatanKeuangan) TableName() string {
	return "catatan_keuangan"
}
