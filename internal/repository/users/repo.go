package users

import (
	"workshop-management/internal/domain"
)

type Users interface {
	Store(m domain.Users) error
	GetByEmail(email string) (domain.Users, error)
}
