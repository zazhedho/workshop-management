package user

type RepoUser interface {
	Store(m Users) error
	GetByEmail(email string) (Users, error)
	GetByID(id string) (Users, error)
	GetAll(page, limit int, orderBy, orderDir, search string) ([]Users, int64, error)
	Update(m Users) error
	Delete(id string) error
}
