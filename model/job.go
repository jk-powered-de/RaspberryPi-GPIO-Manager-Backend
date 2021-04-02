package model

import "github.com/jinzhu/gorm"

type Job struct {
	gorm.Model
	NamedPinId int    `json:"named_pin_id" binding:"required"`
	State      string `json:"state"`
	StartTime  int64  `json:"start_time" binding:"required"`
	EndTime    int64  `json:"end_time" binding:"required"`
}

type JobPatch struct {
	gorm.Model
	NamedPinId int    `json:"named_pin_id"`
	State      string `json:"state"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
}
