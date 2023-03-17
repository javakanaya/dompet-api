package repository

import (
	"context"

	"dompet-api/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if tx := db.connection.Create(&user).Error; tx != nil {
		return entity.User{}, tx
	}

	return user, nil
}
