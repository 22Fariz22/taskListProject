package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type Task struct {
	Id    int
	Time  time.Time
	Title string
	Body  string
}

type AllTasks struct {
	Tasks []*Task
}

func ShowAllTasks() (at *AllTasks) {
	file, err := os.OpenFile("list.json", os.O_RDWR|os.O_APPEND, 0666)
	checkError(err)
	b, err := ioutil.ReadAll(file)
	var alTasks AllTasks
	json.Unmarshal(b, &alTasks.Tasks)
	checkError(err)
	return &alTasks
}
