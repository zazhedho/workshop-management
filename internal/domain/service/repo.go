package service

type RepoService interface {
	Store(m Service) error
	Fetch(page, limit int, orderBy, orderDir, search string) ([]Service, int64, error)
	GetById(id string) (Service, error)
	Update(m Service, data interface{}) (int64, error)
	Delete(m Service, data interface{}) error
}
