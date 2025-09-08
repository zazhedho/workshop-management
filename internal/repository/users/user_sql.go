package users

import (
	"workshop-management/internal/domain"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) Users {
	return &repo{DB: db}
}

func (r *repo) Store(m domain.Users) error {
	return r.DB.Create(&m).Error
}

func (r *repo) GetByEmail(email string) (ret domain.Users, err error) {
	if err = r.DB.Where("email = ?", email).First(&ret).Error; err != nil {
		return domain.Users{}, err
	}

	return ret, nil
}
