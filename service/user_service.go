package service

import (
	"context"
	"dompet-api/dto"
	"dompet-api/entity"
	"dompet-api/repository"
	"dompet-api/utils"

	"github.com/mashingan/smapping"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error)
	IsDuplicateEmail(ctx context.Context, email string) (bool, error)
	VerifyCredential(ctx context.Context, email string, password string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
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
	return us.userRepository.InsertUser(ctx, user)
}

func (us *userService) IsDuplicateEmail(ctx context.Context, email string) (bool, error) {
	checkUser, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if checkUser.Email == "" {
		return false, nil
	}

	return true, nil
}

func (us *userService) VerifyCredential(ctx context.Context, email string, password string) (bool, error) {
	checkUser, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	checkPassword, err := utils.PasswordCompare(checkUser.Password, []byte(password))
	if err != nil {
		return false, err
	}

	if checkUser.Email == email && checkPassword {
		return true, nil
	}
	return false, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	return us.userRepository.GetUserByEmail(ctx, email)
}
