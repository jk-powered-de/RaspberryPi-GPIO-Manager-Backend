package model

import "github.com/jinzhu/gorm"

type NamedPin struct {
	gorm.Model
	Name  string `json:"name" binding:"required"`
	Pin   int    `json:"pin" binding:"required"`
	State string `json:"state"`
}

type NamedPinPatch struct {
	gorm.Model
	Name  string `json:"name"`
	Pin   int    `json:"pin"`
	State string `json:"-"`
}
