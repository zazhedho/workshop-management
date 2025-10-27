package vehicle

import (
	"time"
	"workshop-management/internal/domain/user"

	"gorm.io/gorm"
)

func (Vehicle) TableName() string {
	return "vehicles"
}

type Vehicle struct {
	Id           string         `json:"id" gorm:"column:id;primaryKey"`
	UserId       string         `json:"user_id"`
	LicensePlate string         `json:"license_plate"`
	Brand        string         `json:"brand"`
	Model        string         `json:"model"`
	Year         string         `json:"year"`
	Color        string         `json:"color"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	UpdatedBy    string         `json:"updated_by"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy    string         `json:"-"`

	User user.Users `gorm:"foreignKey:UserId"`
}
