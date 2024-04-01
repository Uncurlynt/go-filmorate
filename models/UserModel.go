package models

type User struct {
	ID       int        `json:"id"`
	Email    string     `json:"email" validate:"required,email"`
	Login    string     `json:"login" validate:"required"`
	Name     string     `json:"name"`
	Birthday CustomTime `json:"birthday"`
}
