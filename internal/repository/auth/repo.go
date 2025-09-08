package auth

import "workshop-management/internal/domain"

type Blacklist interface {
	Store(m domain.Blacklist) error
	GetByToken(token string) (domain.Blacklist, error)
}
