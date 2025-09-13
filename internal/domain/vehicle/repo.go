package vehicle

type RepoVehicle interface {
	Store(m Vehicle) error
	Fetch(page, limit int, orderBy, orderDir, search, userId string) ([]Vehicle, int64, error)
	GetById(id string) (Vehicle, error)
	Update(m Vehicle, data interface{}) (int64, error)
	Delete(m Vehicle, data interface{}) error
}
