package core

import (
	"fmt"
	"strings"
)

type Task struct {
	Id     int
	Name   string
	Status TaskStatus
}

func (t Task) Print() {
	fmt.Printf("%d\t%s\t%s\n", t.Id, strings.TrimSpace(t.Name), t.Status.ToString())
}

func (t *Task) UpdateToInProgress() Task {
	t.Status = InProgress
	return *t
}

func (t *Task) UpdateToCompleted() Task {
	t.Status = Completed

	return *t
}
