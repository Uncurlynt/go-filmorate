package model

import "go-filmorate/types"

type User struct {
	ID       int              `json:"id"`
	Email    string           `json:"email" validate:"required,email"`
	Login    string           `json:"login" validate:"required,validLogin"`
	Name     string           `json:"name"`
	Birthday types.CustomTime `json:"birthday"`
}
