package entity

import (
	"oprec/dompet-api/utils"

	"gorm.io/gorm"
)

type User struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" binding:"required"`
	Email      string    `json:"email" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	ListDompet []*Dompet `gorm:"many2many:detail_user_dompet;" json:"list_dompet,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = utils.PasswordHash(u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	var err error
	if u.Password != "" {
		u.Password, err = utils.PasswordHash(u.Password)
	}
	if err != nil {
		return err
	}
	return nil
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return "users"
}
