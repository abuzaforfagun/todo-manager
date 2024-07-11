package core

type TaskStatus int

const (
	Pending TaskStatus = iota
	InProgress
	Completed
)

func (status TaskStatus) ToString() string {
	return [...]string{"Pending", "In Progress", "Completed", "Canceled"}[status]
}
