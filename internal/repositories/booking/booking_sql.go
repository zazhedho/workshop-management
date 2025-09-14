package booking

import (
	"workshop-management/internal/domain/booking"
	"workshop-management/internal/domain/service"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewBookingRepo(db *gorm.DB) booking.RepoBooking {
	return &repo{DB: db}
}

func (r *repo) Create(booking booking.Booking, bookingServices []booking.BookService) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Omit("Services").Create(&booking).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(bookingServices) > 0 {
		if err := tx.Create(&bookingServices).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *repo) GetServicesByIDs(serviceIDs []string) ([]service.Service, error) {
	var services []service.Service
	if err := r.DB.Where("id IN ?", serviceIDs).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}
