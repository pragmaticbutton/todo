package main

import (
	"fmt"
	"todo/internal/service"
	"todo/internal/storage"
)

func main() {
	data := storage.New()
	svc := service.New(data)

	t, err := svc.AddTask("first task")
	if err != nil {
		panic(err)
	}
	fmt.Println(*t)

	svc.CompleteTask(t.ID)

	tasks, err := svc.ListTasks()
	if err != nil {
		panic(err)
	}
	fmt.Println(tasks)

	svc.ReopenTask(t.ID)

	tasks, err = svc.ListTasks()
	if err != nil {
		panic(err)
	}
	fmt.Println(tasks)

}
