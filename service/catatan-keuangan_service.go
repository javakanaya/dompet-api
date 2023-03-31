package service

import (
	"context"
	"dompet-api/dto"
	"dompet-api/entity"
	"dompet-api/repository"

	"github.com/mashingan/smapping"
)

type catatanService struct {
	catatanRepo repository.CatatanRepository
}

type CatatanService interface {
	CreatePemasukan(ctx context.Context, pemasukanDTO dto.CreatePemasukanDTO) (entity.CatatanKeuangan, error)
	CreatePengeluaran(ctx context.Context, pengeluaranDTO dto.CreatePengeluaranDTO) (entity.CatatanKeuangan, error)
}

func NewCatatanService(cr repository.CatatanRepository) CatatanService {
	return &catatanService{
		catatanRepo: cr,
	}
}

func (s *catatanService) CreatePemasukan(ctx context.Context, pemasukanDTO dto.CreatePemasukanDTO) (entity.CatatanKeuangan, error) {
	var catatanPemasukan entity.CatatanKeuangan
	if err := smapping.FillStruct(&catatanPemasukan, smapping.MapFields(pemasukanDTO)); err != nil {
		return catatanPemasukan, err
	}
	catatanPemasukan.Pengeluaran = 0
	catatanPemasukan.Jenis = "Pemasukan"
	return s.catatanRepo.CreateCatatanKeuangan(ctx, catatanPemasukan)
}

func (s *catatanService) CreatePengeluaran(ctx context.Context, pengeluaranDTO dto.CreatePengeluaranDTO) (entity.CatatanKeuangan, error) {
	var catatanPengeluaran entity.CatatanKeuangan
	if err := smapping.FillStruct(&catatanPengeluaran, smapping.MapFields(pengeluaranDTO)); err != nil {
		return catatanPengeluaran, err
	}
	catatanPengeluaran.Pemasukan = 0
	catatanPengeluaran.Jenis = "Pengeluaran"
	return s.catatanRepo.CreateCatatanKeuangan(ctx, catatanPengeluaran)
}
