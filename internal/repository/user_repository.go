package repository

import (
	"context"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Get(ctx context.Context, userID int) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
