package models

import "go-filmorate/utils"

type User struct {
	ID       int              `json:"id"`
	Email    string           `json:"email" validate:"required,email"`
	Login    string           `json:"login" validate:"required,isValidLogin"`
	Name     string           `json:"name"`
	Birthday utils.CustomTime `json:"birthday"`
	Friends  []int            `json:"friends"`
}
