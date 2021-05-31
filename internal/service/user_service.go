package service

import (
	"github.com/muktiarafi/rodavis-api/internal/entity"
	"github.com/muktiarafi/rodavis-api/internal/model"
)

type UserService interface {
	Create(userDTO *model.CreateUserDTO) (*model.UserDTO, error)
	Auth(loginDTO *model.LoginDTO) (*model.UserDTO, error)
	Get(userID int) (*entity.User, error)
}
