package repository

import (
	"context"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/driver"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, e driver.Executor, user *entity.User) (*entity.User, error)
	Get(ctx context.Context, e driver.Executor, userID int) (*entity.User, error)
	GetByEmail(ctx context.Context, e driver.Executor, email string) (*entity.User, error)
}
