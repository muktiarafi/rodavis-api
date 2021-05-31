package service

import (
	"database/sql"

	"github.com/muktiarafi/rodavis-api/internal/api"
	"github.com/muktiarafi/rodavis-api/internal/entity"
	"github.com/muktiarafi/rodavis-api/internal/model"
	"github.com/muktiarafi/rodavis-api/internal/repository"
	"github.com/muktiarafi/rodavis-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		UserRepository: userRepo,
	}
}

func (s *UserServiceImpl) Create(createUserDTO *model.CreateUserDTO) (*model.UserDTO, error) {
	const op = "UserServiceImpl.Create"
	hash, err := bcrypt.GenerateFromPassword([]byte(createUserDTO.Password), 12)
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"bcrypt.GenerateFromPassword",
			err,
		)
	}

	newUser := &entity.User{
		Name:        createUserDTO.Name,
		Email:       createUserDTO.Email,
		Password:    string(hash),
		PhoneNumber: createUserDTO.PhoneNumber,
	}

	newUser, err = s.UserRepository.Create(newUser)
	if err != nil {
		return nil, err
	}

	token, err := utils.CreateToken(&utils.UserPayload{newUser.ID, newUser.Email, newUser.Role})
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"utils.CreateToken",
			err,
		)
	}

	userDTO := &model.UserDTO{
		User:  newUser,
		Token: token,
	}

	return userDTO, nil
}

func (s *UserServiceImpl) Auth(loginDTO *model.LoginDTO) (*model.UserDTO, error) {
	const op = "UserServiceImpl.Auth"
	user, err := s.UserRepository.GetByEmail(loginDTO.Email)
	if err != nil {
		if exc, ok := err.(*api.Exception); ok {
			if exc.Err == sql.ErrNoRows {
				return nil, api.NewSingleMessageException(
					api.EINVALID,
					op,
					"Invalid Email or Password",
					err,
				)
			}
		}
		return nil, err
	}

	if err := bcrypt.
		CompareHashAndPassword([]byte(user.Password), []byte(loginDTO.Password)); err != nil {
		return nil, api.NewSingleMessageException(
			api.EINVALID,
			op,
			"Invalid Email or Password",
			err,
		)
	}

	token, err := utils.CreateToken(&utils.UserPayload{user.ID, user.Email, user.Role})
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"utils.CreateToken",
			err,
		)
	}

	userDTO := &model.UserDTO{
		User:  user,
		Token: token,
	}

	return userDTO, nil
}

func (s *UserServiceImpl) Get(userID int) (*entity.User, error) {
	return s.UserRepository.Get(userID)
}
