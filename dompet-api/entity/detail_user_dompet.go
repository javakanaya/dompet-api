package entity

type DetailUserDompet struct {
	UserID   uint64 `gorm:"primaryKey" json:"user_id"`
	DompetID uint64 `gorm:"primaryKey" json:"dompet_id"`
}

func (DetailUserDompet) TableName() string {
	return "detail_user_dompet"
}
