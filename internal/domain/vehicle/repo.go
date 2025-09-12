package vehicle

type Repository interface {
	Store(m Vehicle) error
	FetchVehicles(page, limit int, orderBy, orderDir, search, userId string) ([]Vehicle, int64, error)
	GetVehicle(id string) (Vehicle, error)
	UpdateVehicle(id string, data interface{}) (int64, error)
	DeleteVehicle(m Vehicle, data interface{}) error
}
