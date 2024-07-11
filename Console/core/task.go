package core

import "fmt"

type Task struct {
	Id   int
	Name string
}

func (t Task) Print() {
	fmt.Printf("%d \t %s\n", t.Id, t.Name)
}
