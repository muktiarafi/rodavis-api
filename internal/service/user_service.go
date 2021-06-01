package service

import (
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
)

type UserService interface {
	Create(userDTO *model.CreateUserDTO) (*model.UserDTO, error)
	Auth(loginDTO *model.LoginDTO) (*model.UserDTO, error)
	Get(userID int) (*entity.User, error)
}
