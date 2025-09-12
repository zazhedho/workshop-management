package repository

import (
	"fmt"
	"strings"
	"workshop-management/internal/domain/vehicle"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewVehicleRepo(db *gorm.DB) vehicle.Repository {
	return &repo{DB: db}
}

func (r *repo) Store(m vehicle.Vehicle) error {
	return r.DB.Create(&m).Error
}

func (r *repo) FetchVehicles(page, limit int, orderBy, orderDir, search, userId string) (ret []vehicle.Vehicle, totalData int64, err error) {
	query := r.DB.Table(vehicle.Vehicle{}.TableName())

	if userId != "" {
		query = query.Where("user_id = ?", userId)
	}
	if strings.TrimSpace(search) != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("LOWER(license_plate) LIKE LOWER(?) OR LOWER(brand) LIKE LOWER(?) OR LOWER(model) LIKE LOWER(?) OR year LIKE ? OR LOWER(color) LIKE LOWER(?)", searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	if orderBy != "" && orderDir != "" {
		validColumns := map[string]bool{
			"license_plate": true,
			"brand":         true,
			"model":         true,
			"year":          true,
			"created_at":    true,
			"updated_at":    true,
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

func (r *repo) GetVehicle(id string) (vehicle.Vehicle, error) {
	var m vehicle.Vehicle
	if err := r.DB.Where("id = ?", id).First(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (r *repo) UpdateVehicleById(id string, m vehicle.Vehicle) error {
	//TODO implement me
	panic("implement me")
}

func (r *repo) DeleteVehicleById(id string) error {
	//TODO implement me
	panic("implement me")
}
