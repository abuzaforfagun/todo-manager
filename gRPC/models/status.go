package models

type Status int

const (
	Pending Status = iota
	InProgress
	Completed
)

func (status Status) ToString() string {
	return [...]string{"Pending", "In Progress", "Completed", "Canceled"}[status]
}
