package service

import (
	"context"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
)

type UserService interface {
	Create(ctx context.Context, userDTO *model.CreateUserDTO) (*model.UserDTO, error)
	Auth(ctx context.Context, loginDTO *model.LoginDTO) (*model.UserDTO, error)
	Get(ctx context.Context, userID int) (*entity.User, error)
}
