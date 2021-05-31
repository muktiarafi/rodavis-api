package repository

import "github.com/muktiarafi/rodavis-api/internal/entity"

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	Get(userID int) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
}
