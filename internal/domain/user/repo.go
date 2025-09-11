package user

type Repository interface {
	Store(m Users) error
	GetByEmail(email string) (Users, error)
	GetByID(id string) (Users, error)
	GetAll() ([]Users, error)
	Update(m Users) error
	Delete(id string) error
}
