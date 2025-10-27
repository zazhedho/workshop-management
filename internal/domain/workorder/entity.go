package workorder

import "time"

type WorkOrder struct {
	Id         string `json:"id" gorm:"type:uuid;primaryKey"`
	BookingId  string `json:"booking_id" gorm:"type:uuid;not null"`
	CustomerId string `json:"customer_id" gorm:"type:uuid;not null"`
	VehicleId  string `json:"vehicle_id" gorm:"type:uuid;not null"`
	MechanicId string `json:"mechanic_id"`
	Status     string `json:"status"` // pending, in_progress, completed
	Notes      string `json:"notes"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	DeletedAt time.Time `json:"-"`
	DeletedBy string    `json:"-"`

	Services []SvcWorkOrder  `gorm:"foreignKey:WorkOrderId" json:"services"`
	Parts    []PartWorkOrder `gorm:"foreignKey:WorkOrderId" json:"parts"`
}

type PartWorkOrder struct {
	Id          string    `json:"id"`
	WorkOrderId string    `json:"work_order_id"`
	SparepartId string    `json:"sparepart_id"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

type SvcWorkOrder struct {
	Id          string  `json:"id"`
	WorkOrderId string  `json:"work_order_id"`
	ServiceId   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Status      string  `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	DeletedAt time.Time `json:"-"`
	DeletedBy string    `json:"-"`
}
