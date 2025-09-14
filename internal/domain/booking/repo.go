package booking

import "workshop-management/internal/domain/service"

type RepoBooking interface {
	Create(booking Booking, bookingServices []BookService) error
	GetServicesByIDs(serviceIDs []string) ([]service.Service, error)
}
