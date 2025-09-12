package repository

import (
	"fmt"
	"strings"
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

func (r *repo) GetAll(page, limit int, orderBy, orderDir, search string) (ret []user.Users, totalData int64, err error) {
	query := r.DB.Table(user.Users{}.TableName()).Debug()

	if strings.TrimSpace(search) != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?) OR LOWER(phone) LIKE LOWER(?) OR LOWER(role) LIKE LOWER(?)", searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if err := query.Where("deleted_at IS NULL").Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	if orderBy != "" && orderDir != "" {
		validColumns := map[string]bool{
			"name":       true,
			"email":      true,
			"phone":      true,
			"role":       true,
			"created_at": true,
			"updated_at": true,
		}

		validDirections := map[string]bool{
			"asc":  true,
			"desc": true,
		}

		if _, ok := validColumns[orderBy]; !ok {
			return nil, 0, fmt.Errorf("invalid orderBy column: %s", orderBy)
		}
		if _, ok := validDirections[orderDir]; !ok {
			return nil, 0, fmt.Errorf("invalid orderDir: %s", orderDir)
		}

		query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDir))
	}

	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if err = query.Find(&ret).Error; err != nil {
		return ret, 0, err
	}

	return ret, totalData, nil
}

func (r *repo) Update(m user.Users) error {
	return r.DB.Save(&m).Error
}

func (r *repo) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&user.Users{}).Error
}
