package auth

type Repository interface {
	Store(m Blacklist) error
	GetByToken(token string) (Blacklist, error)
}
