package model

import "go-filmorate/types"

type Film struct {
	ID          int              `json:"id"`
	Name        string           `json:"name" validate:"required"`
	Description string           `json:"description" validate:"required,lte=200"`
	ReleaseDate types.CustomTime `json:"releaseDate" validate:"required"`
	Duration    int              `json:"duration" validate:"required,gt=0"`
}
