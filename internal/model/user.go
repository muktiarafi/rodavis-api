package model

import "github.com/muktiarafi/rodavis-api/internal/entity"

type CreateUserDTO struct {
	Name        string `json:"name" validate:"required,min=4"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=15"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
}

type UserDTO struct {
	User  *entity.User `json:"user"`
	Token string       `json:"token"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUser struct {
	Name string `validate:"min=4"`
}
