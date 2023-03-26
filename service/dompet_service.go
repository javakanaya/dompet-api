package service

import (
	"context"
	"dompet-api/dto"
	"dompet-api/entity"
	"dompet-api/repository"

	"github.com/mashingan/smapping"
)

type dompetService struct {
	dompetRepo repository.DompetRepository
}

type DompetService interface {
	GetMyDompet(id uint64) (entity.User, error)
	CreateDompet(ctx context.Context, dompetDTO dto.DompetCreateDTO) (entity.Dompet, error)
}

func NewDompetService(dr repository.DompetRepository) DompetService {
	return &dompetService{
		dompetRepo: dr,
	}
}

func (s *dompetService) GetMyDompet(id uint64) (entity.User, error) {
	berhasilGet, err := s.dompetRepo.GetMyDompet(nil, id)
	if err != nil {
		return entity.User{}, err
	}

	return berhasilGet, nil
}

func (s *dompetService) CreateDompet(ctx context.Context, dompetDTO dto.DompetCreateDTO) (entity.Dompet, error) {
	var dompet entity.Dompet
	if err := smapping.FillStruct(&dompet, smapping.MapFields(&dompetDTO)); err != nil {
		return dompet, err
	}
	return s.dompetRepo.InsertDompet(ctx, dompet)
}
