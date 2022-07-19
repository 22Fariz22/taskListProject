package main

import "time"

type Task struct {
	id    int
	time  time.Time
	title string
}
type AllTasks struct {
	Tasks []*Task
}

func ShowAllTasks() (at *AllTasks) {
	return
}
