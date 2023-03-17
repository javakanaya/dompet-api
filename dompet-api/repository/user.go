package repository

import (
	"context"
	"oprec/dompet-api/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	// functional
	CreateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	var err error
	if tx == nil {
		r.db.WithContext(ctx).Debug().Create(&user)
	} else {
		err = tx.WithContext(ctx).Debug().Create(&user).Error
	}

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error) {
	var err error
	var user entity.User

	if tx == nil {
		r.db.WithContext(ctx).Debug().Where(("email = ?"), email).Take(&user)
	} else {
		err = tx.WithContext(ctx).Debug().Where(("email = ?"), email).Take(&user).Error
	}

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
