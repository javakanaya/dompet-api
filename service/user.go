package service

import (
	"context"

	"dompet-api/dto"
	"dompet-api/entity"
	"dompet-api/repository"

	"github.com/mashingan/smapping"
)

type UserService interface {
	CreateUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func (us *userService) CreateUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error) {
	var user entity.User
	if err := smapping.FillStruct(&user, smapping.MapFields(&userDTO)); err != nil {
		return user, err
	}

	return us.userRepository.CreateUser(ctx, user)
}
