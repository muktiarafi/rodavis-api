package service

import (
	"context"
	"database/sql"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/api"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/repository"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/utils"
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

func (s *UserServiceImpl) Create(ctx context.Context, createUserDTO *model.CreateUserDTO) (*model.UserDTO, error) {
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

	newUser, err = s.UserRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	token, err := utils.CreateToken(&model.UserPayload{newUser.ID, newUser.Email, newUser.Role})
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

func (s *UserServiceImpl) Auth(ctx context.Context, loginDTO *model.LoginDTO) (*model.UserDTO, error) {
	const op = "UserServiceImpl.Auth"
	user, err := s.UserRepository.GetByEmail(ctx, loginDTO.Email)
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

	token, err := utils.CreateToken(&model.UserPayload{user.ID, user.Email, user.Role})
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

func (s *UserServiceImpl) Get(ctx context.Context, userID int) (*entity.User, error) {
	return s.UserRepository.Get(ctx, userID)
}
