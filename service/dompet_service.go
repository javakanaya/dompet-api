package service

import (
	"dompet-api/entity"
	"dompet-api/repository"
)

type dompetService struct {
	dompetRepo repository.DompetRepository
}

type DompetService interface {
	GetMyDompet(id uint64) (entity.User, error)
	GetDetailDompet(id uint64) (entity.Dompet, error)
	InviteToDompet(inviteDTO dto.InviteUserRequest) (entity.User, error)
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

func (s *dompetService) GetDetailDompet(id uint64) (entity.Dompet, error) {
	berhasilGet, err := s.dompetRepo.GetDetailDompet(nil, id)
	if err != nil {
		return entity.Dompet{}, err
	}

	return berhasilGet, nil
}

func (s *dompetService) InviteToDompet(inviteDTO dto.InviteUserRequest) (entity.User, error) {
	berhasilInvite, err := s.dompetRepo.InviteToDompet(nil, inviteDTO.DompetID, inviteDTO.EmailUser)
	if err != nil {
		return entity.User{}, err
	}

	return berhasilInvite, nil
}
