package repository

import (
	"context"
	"dompet-api/dto"
	"dompet-api/entity"
	"errors"

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
	CreateCatatanKeuangan(ctx context.Context, catatanKeuangan entity.CatatanKeuangan) (entity.CatatanKeuangan, error)
	DeleteCatatanKeuangan(ctx context.Context, catatanKeuanganID uint64) error
	GetCatatanByID(ctx context.Context, catatanKeuanganID uint64) (entity.CatatanKeuangan, error)
	UpdateCatatan(ctx context.Context, catatanKeuangan entity.CatatanKeuangan) (entity.CatatanKeuangan, error)
	GetKategori(ctx context.Context, KategoriCatatanKeuangan string) (entity.KategoriCatatanKeuangan, error)
	GetListKategori(jenis string) ([]dto.ReturnKategori, error)
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

func (r *catatanRepository) CreateCatatanKeuangan(ctx context.Context, catatanPemasukan entity.CatatanKeuangan) (entity.CatatanKeuangan, error) {
	if tx := r.db.Create(&catatanPemasukan).Error; tx != nil {
		return entity.CatatanKeuangan{}, tx
	}

	return catatanPemasukan, nil
}

func (r *catatanRepository) DeleteCatatanKeuangan(ctx context.Context, catatanKeuanganID uint64) error {
	if tx := r.db.Delete(&entity.CatatanKeuangan{}, "id = ?", &catatanKeuanganID).Error; tx != nil {
		return tx
	}
	return nil
}
func (r *catatanRepository) GetCatatanByID(ctx context.Context, catatanKeuanganID uint64) (entity.CatatanKeuangan, error) {
	var catatanKeuangan entity.CatatanKeuangan
	if tx := r.db.Where("id = ?", catatanKeuanganID).Take(&catatanKeuangan).Error; tx != nil {
		return entity.CatatanKeuangan{}, tx
	}
	return catatanKeuangan, nil
}

func (r *catatanRepository) UpdateCatatan(ctx context.Context, catatanKeuangan entity.CatatanKeuangan) (entity.CatatanKeuangan, error) {
	if tx := r.db.Save(&catatanKeuangan).Error; tx != nil {
		return entity.CatatanKeuangan{}, tx
	}
	return catatanKeuangan, nil
}

func (r *catatanRepository) GetKategori(ctx context.Context, KategoriCatatanKeuangan string) (entity.KategoriCatatanKeuangan, error) {
	var kategori entity.KategoriCatatanKeuangan
	if tx := r.db.Where("nama_kategori = ?", KategoriCatatanKeuangan).Take(&kategori).Error; tx != nil {
		return entity.KategoriCatatanKeuangan{}, tx
	}
	return kategori, nil
}
func (r *catatanRepository) GetListKategori(jenis string) ([]dto.ReturnKategori, error) {
	var model entity.KategoriCatatanKeuangan
	var ListKategori []dto.ReturnKategori
	var getKategori *gorm.DB

	if jenis == "pemasukan" {
		getKategori = r.db.Debug().Model(&model).Where("jenis_id = 1").Find(&ListKategori)
	} else if jenis == "pengeluaran" {
		getKategori = r.db.Debug().Model(&model).Where("jenis_id = 2").Find(&ListKategori)
	} else {
		return []dto.ReturnKategori{}, errors.New("invalid jenis kategori")
	}

	if getKategori.Error != nil {
		return []dto.ReturnKategori{}, getKategori.Error
	}
	return ListKategori, nil

}
