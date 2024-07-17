package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name   string
	Status Status
	UserId uint
	User   User
}

func (Task) TableName() string {
	return "Tasks"
}

type TaskDto struct {
	Id        uint
	Name      string
	Status    string
	CreatedAt time.Time
}

type TaskRequestDto struct {
	Name string
}
