package workorder

import (
	"workshop-management/internal/domain/workorder"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewWorkOrderRepo(db *gorm.DB) workorder.RepoWorkOrder {
	return &repo{DB: db}
}

func (r *repo) Create(workOrder workorder.WorkOrder, svcWorkOrders []workorder.SvcWorkOrder) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Omit("Services").Create(&workOrder).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(svcWorkOrders) > 0 {
		if err := tx.Create(&svcWorkOrders).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *repo) GetById(id string) (workorder.WorkOrder, error) {
	var wo workorder.WorkOrder
	if err := r.DB.Where("id = ?", id).First(&wo).Error; err != nil {
		return workorder.WorkOrder{}, err
	}

	return wo, nil
}

func (r *repo) Update(m workorder.WorkOrder, data map[string]interface{}) (int64, error) {
	res := r.DB.Model(&m).Where("id = ?", m.Id).Updates(data)
	if res.Error != nil {
		return 0, res.Error
	}

	return res.RowsAffected, nil
}
