package main

import (
	"todo/internal/service"
	"todo/internal/storage/memory"
)

func main() {
	svc := service.New(memory.New())
	svc.AddTask("Task 1")
	svc.AddTask("Task 2")
	tasks, _ := svc.ListTasks()
	for _, t := range tasks {
		println(t.ID, t.Description, t.Done)
	}
	svc.CompleteTask(2)
	tasks, _ = svc.ListTasks()
	for _, t := range tasks {
		println(t.ID, t.Description, t.Done)
	}

}
