package booking

import "time"

type Booking struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	VehicleID string    `json:"vehicle_id"`
	ServiceID string    `json:"service_id"`
	Status    string    `json:"status"` // pending, confirmed, cancelled
	BookingAt time.Time `json:"booking_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
