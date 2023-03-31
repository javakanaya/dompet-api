package service

import (
	"dompet-api/dto"
	"dompet-api/entity"
	"dompet-api/repository"
)

type catatanService struct {
	catatanRepo repository.CatatanRepository
}

type CatatanService interface {
	Transfer(transferDTO dto.TransferRequest, idUser uint64, idSumber uint64) (entity.CatatanKeuangan, error)
	InsertKategori(kategori entity.KategoriCatatanKeuangan) (entity.KategoriCatatanKeuangan, error)
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
