package entity

type KategoriCatatanKeuangan struct {
	// ini jenis nya gmn enaknya string atau uint?
	JenisID      uint   `gorm:"primaryKey;autoIncrement:false" json:"jenis_id"`
	ID           uint64 `gorm:"primaryKey" json:"id"`
	NamaKategori string `json:"nama_kategori"`
}
