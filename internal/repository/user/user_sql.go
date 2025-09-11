package repository

import (
	"workshop-management/internal/domain/user"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) user.Repository {
	return &repo{DB: db}
}

func (r *repo) Store(m user.Users) error {
	return r.DB.Create(&m).Error
}

func (r *repo) GetByEmail(email string) (ret user.Users, err error) {
	if err = r.DB.Where("email = ?", email).First(&ret).Error; err != nil {
		return user.Users{}, err
	}

	return ret, nil
}

func (r *repo) GetByID(id string) (ret user.Users, err error) {
	if err = r.DB.Where("id = ?", id).First(&ret).Error; err != nil {
		return user.Users{}, err
	}
	return ret, nil
}

func (r *repo) GetAll() (ret []user.Users, err error) {
	if err = r.DB.Find(&ret).Error; err != nil {
		return []user.Users{}, err
	}
	return ret, nil
}

func (r *repo) Update(m user.Users) error {
	return r.DB.Save(&m).Error
}

func (r *repo) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&user.Users{}).Error
}
