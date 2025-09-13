package service

import (
	"fmt"
	"strings"
	"workshop-management/internal/domain/service"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewServiceRepo(db *gorm.DB) service.RepoService {
	return &repo{DB: db}
}

func (r *repo) Store(m service.Service) error {
	return r.DB.Create(&m).Error
}

func (r *repo) Fetch(page, limit int, orderBy, orderDir, search string) (ret []service.Service, totalData int64, err error) {
	query := r.DB.Table(service.Service{}.TableName()).Debug()

	if strings.TrimSpace(search) != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", searchPattern, searchPattern)
	}

	if err := query.Where("deleted_at IS NULL").Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	if orderBy != "" && orderDir != "" {
		validColumns := map[string]bool{
			"name":       true,
			"price":      true,
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

func (r *repo) GetById(id string) (service.Service, error) {
	var m service.Service
	if err := r.DB.Where("id = ?", id).First(&m).Error; err != nil {
		return service.Service{}, err
	}
	return m, nil
}

func (r *repo) Update(m service.Service, data interface{}) (int64, error) {
	res := r.DB.Table(m.TableName()).Where("id = ?", m.Id).Updates(data)
	if res.Error != nil {
		return 0, nil
	}
	return res.RowsAffected, nil
}

func (r *repo) Delete(m service.Service, data interface{}) error {
	return r.DB.Model(&m).Where("id = ?", m.Id).Updates(data).Error
}
