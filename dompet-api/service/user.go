package service

import (
	"context"
	"errors"
	"oprec/dompet-api/dto"
	"oprec/dompet-api/entity"
	"oprec/dompet-api/repository"

	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
)

type userService struct {
	userRepo repository.UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, userDTO dto.UserRegisterRequest) (entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepo: ur,
	}
}

func (s *userService) CreateUser(ctx context.Context, userDTO dto.UserRegisterRequest) (entity.User, error) {
	var user entity.User
	copier.Copy(&user, &userDTO)

	checkEmail, err := s.userRepo.FindUserByEmail(ctx, nil, user.Email)
	if err != nil {
		return entity.User{}, err
	}

	if !(cmp.Equal(checkEmail, entity.User{})) { // saya menggunakan library cmp dengan tujuan untuk membandingkan 2 struct, tidak bisa dengan = karena pada struct user terdapat []blog
		return entity.User{}, errors.New("email yang diinput sudah pernah digunakan")
	}

	berhasilRegis, err := s.userRepo.CreateUser(ctx, nil, user)
	if err != nil {
		return entity.User{}, err
	}

	return berhasilRegis, nil

}

func (s *userService) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, nil, email)
	if err != nil {
		return entity.User{}, err
	}

	if cmp.Equal(user, entity.User{}) {
		return entity.User{}, errors.New("email tidak valid")
	}

	return user, nil
}
