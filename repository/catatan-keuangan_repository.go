package repository

import (
	"errors"
	"dompet-api/entity"
	"time"

	"gorm.io/gorm"
)

type catatanRepository struct {
	db *gorm.DB
}

type CatatanRepository interface {
	// functional
	Transfer(tx *gorm.DB, idUser uint64, idSumber uint64, tujuan string, nominal uint64, deskripsi string, kategori string) (entity.CatatanKeuangan, error)
	InsertKategori(kategori entity.KategoriCatatanKeuangan) (entity.KategoriCatatanKeuangan, error)
}

func NewCatatanRepository(db *gorm.DB) CatatanRepository {
	return &catatanRepository{
		db: db,
	}
}

func (r *catatanRepository) Transfer(tx *gorm.DB, idUser uint64, idSumber uint64, tujuan string, nominal uint64, deskripsi string, kategori string) (entity.CatatanKeuangan, error) {
	var dompetSumber entity.Dompet
	var dompetTujuan entity.Dompet
	var kategoriTransfer entity.KategoriCatatanKeuangan

	getSumber := r.db.Where("user_id = ? AND id = ?", idUser, idSumber).Take(&dompetSumber)
	if getSumber.Error != nil {
		return entity.CatatanKeuangan{}, errors.New("dompet sumber tidak ditemukan")
	}

	if dompetSumber.Saldo < nominal {
		return entity.CatatanKeuangan{}, errors.New("saldo tidak mencukupi")
	}

	getTujuan := r.db.Where("nama_dompet = ?", tujuan).Take(&dompetTujuan)
	if getTujuan.Error != nil {
		return entity.CatatanKeuangan{}, errors.New("dompet tujuan tidak ditemukan")
	}

	getKategori := r.db.Where("nama_kategori = ?", kategori).Take(&kategoriTransfer)
	if getKategori.Error != nil {
		return entity.CatatanKeuangan{}, errors.New("kategori tidak ditemukan")
	}

	// update saldo sumber
	r.db.Debug().Model(&dompetSumber).Where(entity.Dompet{ID: dompetSumber.ID}).Update("saldo", (dompetSumber.Saldo - nominal))

	// buat catatan di sumber
	newCatatanSumber := entity.CatatanKeuangan{
		Deskripsi:   deskripsi,
		Pemasukan:   0,
		Pengeluaran: nominal,
		Tanggal:     time.Now(),

		Jenis:    "Transfer",
		Kategori: kategoriTransfer.NamaKategori,

		DompetID: dompetSumber.ID,
	}

	if err := r.db.Create(&newCatatanSumber).Error; err != nil {
		return entity.CatatanKeuangan{}, errors.New("gagal membuat catatan keuangan")
	}

	// update saldo tujuan
	r.db.Debug().Model(&dompetTujuan).Where(entity.Dompet{ID: dompetTujuan.ID}).Update("saldo", (dompetTujuan.Saldo + nominal))

	// buat catatan di tujuan
	sumber := dompetSumber.NamaDompet
	newCatatanTujuan := entity.CatatanKeuangan{
		Deskripsi:   deskripsi,
		Pemasukan:   nominal,
		Pengeluaran: 0,
		Tanggal:     time.Now(),

		Jenis:    "Pemasukan",
		Kategori: "Menerima transfer dari dompet: " + sumber,

		DompetID: dompetTujuan.ID,
	}

	if err := r.db.Create(&newCatatanTujuan).Error; err != nil {
		return entity.CatatanKeuangan{}, errors.New("gagal membuat catatan keuangan")
	}

	return newCatatanSumber, nil
}

func (r *catatanRepository) InsertKategori(kategori entity.KategoriCatatanKeuangan) (entity.KategoriCatatanKeuangan, error) {
	if err := r.db.Create(&kategori).Error; err != nil {
		return entity.KategoriCatatanKeuangan{}, errors.New("gagal insert kategori")
	}

	return kategori, nil
}
