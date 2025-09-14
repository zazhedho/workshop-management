package booking

import (
	"time"
	"workshop-management/internal/domain/booking"
	"workshop-management/internal/dto"
	"workshop-management/utils"
)

type ServiceBooking struct {
	BookingRepo booking.RepoBooking
}

func NewServiceBooking(bookingRepo booking.RepoBooking) *ServiceBooking {
	return &ServiceBooking{
		BookingRepo: bookingRepo,
	}
}

func (s *ServiceBooking) Create(userId string, req dto.CreateBooking) (booking.Booking, error) {
	bookingID := utils.CreateUUID()
	bookingData := booking.Booking{
		Id:          bookingID,
		UserId:      userId,
		VehicleId:   req.VehicleID,
		BookingDate: req.BookingDate,
		Notes:       req.Notes,
		Status:      utils.StsPending,
		CreatedAt:   time.Now(),
	}

	var bookingServices []booking.BookService

	for _, serviceID := range req.ServiceIDs {
		bookingServices = append(bookingServices, booking.BookService{
			Id:        utils.CreateUUID(),
			BookingID: bookingID,
			ServiceID: serviceID,
		})
	}
	dataService, err := s.BookingRepo.GetServicesByIDs(req.ServiceIDs)
	if err != nil {
		return booking.Booking{}, err
	}
	bookingData.Services = dataService

	if err := s.BookingRepo.Create(bookingData, bookingServices); err != nil {
		return booking.Booking{}, err
	}

	return bookingData, nil
}
