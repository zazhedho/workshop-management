package auth

import (
	"workshop-management/internal/domain"

	"gorm.io/gorm"
)

type blacklistRepo struct {
	DB *gorm.DB
}

func NewBlacklistRepo(db *gorm.DB) Blacklist {
	return &blacklistRepo{
		DB: db,
	}
}

func (r *blacklistRepo) Store(blacklist domain.Blacklist) error {
	return r.DB.Create(&blacklist).Error
}

func (r *blacklistRepo) GetByToken(token string) (domain.Blacklist, error) {
	var blacklist domain.Blacklist
	err := r.DB.Where("token = ?", token).First(&blacklist).Error
	return blacklist, err
}
