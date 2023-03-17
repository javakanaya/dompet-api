package service

import (
	"oprec/dompet-api/entity"
	"oprec/dompet-api/repository"
)

type dompetService struct {
	dompetRepo repository.DompetRepository
}

type DompetService interface {
	GetMyDompet(id uint64) (entity.User, error)
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
