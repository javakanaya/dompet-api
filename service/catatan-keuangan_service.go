package service

import (
	"context"
	"dompet-api/dto"
	"dompet-api/entity"
	"dompet-api/repository"
	"time"

	"github.com/mashingan/smapping"
)

type catatanService struct {
	catatanRepo repository.CatatanRepository
}

type CatatanService interface {
	Transfer(transferDTO dto.TransferRequest, idUser uint64, idSumber uint64) (entity.CatatanKeuangan, error)
	InsertKategori(kategori entity.KategoriCatatanKeuangan) (entity.KategoriCatatanKeuangan, error)
	CreatePemasukan(ctx context.Context, pemasukanDTO dto.CreatePemasukanDTO) (entity.CatatanKeuangan, error)
	CreatePengeluaran(ctx context.Context, pengeluaranDTO dto.CreatePengeluaranDTO) (entity.CatatanKeuangan, error)
	DeleteCatatanKeuangan(ctx context.Context, catatanKeuanganID uint64) error
	IsCatatanExistInDompet(ctx context.Context, catatanKeuanganID uint64, dompetID uint64) (bool, error)
	GetCatatanByID(ctx context.Context, catatanKeuanganId uint64) (entity.CatatanKeuangan, error)
}

func NewCatatanService(cr repository.CatatanRepository) CatatanService {
	return &catatanService{
		catatanRepo: cr,
	}
}

func (s *catatanService) Transfer(transferDTO dto.TransferRequest, idUser uint64, idSumber uint64) (entity.CatatanKeuangan, error) {
	berhasilTransfer, err := s.catatanRepo.Transfer(nil, idUser, idSumber, transferDTO.NamaDompet, transferDTO.Nominal, transferDTO.Deskripsi, transferDTO.Kategori)
	if err != nil {
		return entity.CatatanKeuangan{}, err
	}

	return berhasilTransfer, nil
}

func (s *catatanService) InsertKategori(kategori entity.KategoriCatatanKeuangan) (entity.KategoriCatatanKeuangan, error) {
	berhasilInsert, err := s.catatanRepo.InsertKategori(kategori)
	if err != nil {
		return entity.KategoriCatatanKeuangan{}, err
	}

	return berhasilInsert, nil
}

func (s *catatanService) CreatePemasukan(ctx context.Context, pemasukanDTO dto.CreatePemasukanDTO) (entity.CatatanKeuangan, error) {
	var catatanPemasukan entity.CatatanKeuangan
	if err := smapping.FillStruct(&catatanPemasukan, smapping.MapFields(pemasukanDTO)); err != nil {
		return catatanPemasukan, err
	}
	catatanPemasukan.Tanggal = time.Now()
	catatanPemasukan.Pengeluaran = 0
	catatanPemasukan.Jenis = "Pemasukan"
	return s.catatanRepo.CreateCatatanKeuangan(ctx, catatanPemasukan)
}

func (s *catatanService) CreatePengeluaran(ctx context.Context, pengeluaranDTO dto.CreatePengeluaranDTO) (entity.CatatanKeuangan, error) {
	var catatanPengeluaran entity.CatatanKeuangan
	if err := smapping.FillStruct(&catatanPengeluaran, smapping.MapFields(pengeluaranDTO)); err != nil {
		return catatanPengeluaran, err
	}
	catatanPengeluaran.Tanggal = time.Now()
	catatanPengeluaran.Pemasukan = 0
	catatanPengeluaran.Jenis = "Pengeluaran"
	return s.catatanRepo.CreateCatatanKeuangan(ctx, catatanPengeluaran)
}

func (s *catatanService) DeleteCatatanKeuangan(ctx context.Context, catatanKeuanganID uint64) error {
	return s.catatanRepo.DeleteCatatanKeuangan(ctx, catatanKeuanganID)
}

func (s *catatanService) IsCatatanExistInDompet(ctx context.Context, catatanKeuanganID uint64, dompetID uint64) (bool, error) {
	catatan, err := s.catatanRepo.GetCatatanByID(ctx, catatanKeuanganID)
	if err != nil {
		return false, err
	}
	if catatan.DompetID == dompetID {
		return true, nil
	}
	return false, nil
}

func (s *catatanService) GetCatatanByID(ctx context.Context, catatanKeuanganId uint64) (entity.CatatanKeuangan, error) {
	return s.catatanRepo.GetCatatanByID(ctx, catatanKeuanganId)
}

func (s *catatanService) GetKategori(jenis string) ([]dto.ReturnKategori, error) {
	return s.catatanRepo.GetKategori(jenis)
}
