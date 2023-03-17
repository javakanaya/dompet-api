package entity

type Dompet struct {
	ID                  uint64            `json:"id" gorm:"primaryKey"`
	NamaDompet          string            `json:"nama_dompet" binding:"required"` // keterangan tujuan dibuatnya dompet
	Saldo               uint64            `json:"saldo" binding:"required"`
	ListCatatanKeuangan []CatatanKeuangan `json:"list_catatan_keuangan,omitempty"`
	UserID              uint64            `gorm:"foreignKey" json:"user_id"`
	ListUser            []*User           `gorm:"many2many:detail_user_dompet;" json:"list_user,omitempty"`
}

func (Dompet) TableName() string {
	return "dompet"
}
