package models

import "time"

type Task struct {
	Id        int
	Name      string
	Status    Status
	CreatedAt time.Time
}
