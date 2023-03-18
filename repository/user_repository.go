package repository

import (
	"context"
	"oprec/dompet-api/entity"

	"gorm.io/gorm"
)

type userConnection struct {
	connection *gorm.DB
}

type UserRepository interface {
	// functional
	InsertUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, userID uint64) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := db.connection.Create(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (db *userConnection) GetUserByID(ctx context.Context, userID uint64) (entity.User, error) {
	var user entity.User
	if err := db.connection.Where("id = ?", userID).Take(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (db *userConnection) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	if err := db.connection.Where("email = ?", email).Take(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}
