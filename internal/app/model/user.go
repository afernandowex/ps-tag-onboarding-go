package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;"`
	FirstName string    `gorm:"not null" json:"firstName"`
	LastName  string    `gorm:"not null" json:"lastName"`
	Email     string    `gorm:"not null" json:"email"`
	Age       int8      `gorm:"not null" json:"age"`
}

/*
 * BeforeCreate is a gorm hook that will be called before creating a new user
 * see https://gorm.io/docs/hooks.html for more info
 */
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}
