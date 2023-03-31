package entity

import "time"

type CatatanKeuangan struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	Deskripsi   string    `json:"deskripsi" binding:"required"`   // berupa detail dari pemasukan/pengeluaran
	Pemasukan   uint64    `json:"pemasukan" binding:"required"`   // tuliskan nominal pemasukan
	Pengeluaran uint64    `json:"pengeluaran" binding:"required"` // tuliskan nominal pengeluaran
	Tanggal     time.Time `json:"tanggal" binding:"required"`

	Jenis    string `json:"jenis"`
	Kategori string `json:"kategori"`

	DompetID uint64  `gorm:"foreignKey" json:"dompet_id"`
	Dompet   *Dompet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}

func (CatatanKeuangan) TableName() string {
	return "catatan_keuangan"
}
