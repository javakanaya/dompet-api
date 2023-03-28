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
	GetDetailDompet(id uint64) (entity.Dompet, error)
	InviteToDompet(inviteDTO dto.InviteUserRequest) (entity.User, error)
	IsDompetOwnedByUserID(ctx context.Context, dompetID uint64, userID uint64) (bool, error)
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

// ini verifikasi nya
func (s *dompetService) IsDompetOwnedByUserID(ctx context.Context, dompetID uint64, userID uint64) (bool, error) {
	checkID, err := s.dompetRepo.GetUserIDFromDompet(ctx, dompetID)
	if err != nil {
		return false, err
	}
	if checkID == userID {
		return true, nil
	}
	return false, nil
}
