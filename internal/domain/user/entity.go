package user

import (
	"time"

	"gorm.io/gorm"
)

func (Users) TableName() string {
	return "users"
}

type Users struct {
	Id        string         `json:"id" gorm:"column:id;primaryKey"`
	Name      string         `json:"name" gorm:"column:name"`
	Email     string         `json:"email" gorm:"column:email"`
	Phone     string         `json:"phone" gorm:"column:phone"`
	Password  string         `json:"-" gorm:"column:password"`
	Role      string         `json:"role" gorm:"column:role"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
