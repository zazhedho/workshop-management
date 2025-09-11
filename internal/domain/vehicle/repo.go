package vehicle

type Repository interface {
	Store(m Vehicle) error
	FetchVehicles(page, limit int, orderBy, orderDir, search, userId string) ([]Vehicle, int64, error)
	GetVehicle(id string) (Vehicle, error)
	UpdateVehicleById(id string, m Vehicle) error
	DeleteVehicleById(id string) error
}
