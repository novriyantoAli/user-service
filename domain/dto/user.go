package dto

import "github.com/google/uuid"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	UUID  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role,omitempty"`
	Phone string    `json:"phone"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type RegisterRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
	RoleID          uint
}

type RegisterRespose struct {
	User UserResponse `json:"user"`
}

type UpdateRequest struct {
	Name            string  `json:"name" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	Phone           string  `json:"phone" validate:"required"`
	Password        *string `json:"password,omitempty"`
	ConfirmPassword *string `json:"confirmPassword,omitempty"`
	RoleID          uint
}

type UpdateResponse struct {
	User UserResponse `json:"user"`
}
