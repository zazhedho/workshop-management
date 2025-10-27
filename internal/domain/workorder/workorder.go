package workorder

import "time"

type WorkOrder struct {
	ID         string    `json:"id"`
	BookingID  string    `json:"booking_id"`
	MechanicID string    `json:"mechanic_id"`
	Status     string    `json:"status"` // pending, in_progress, completed
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type WorkOrderPart struct {
	ID          string    `json:"id"`
	WorkOrderID string    `json:"work_order_id"`
	SparepartID string    `json:"sparepart_id"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}
