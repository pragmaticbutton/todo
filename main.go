package main

import (
	"fmt"
	"todo/internal/service"
	"todo/internal/storage/memory"
)

func main() {
	svc := service.New(memory.New())
	svc.AddTask("Task 1")
	svc.AddTask("Task 2")
	tasks, _ := svc.ListTasks()
	for _, t := range tasks {
		fmt.Println(t)
	}
	svc.CompleteTask(2)
	tasks, _ = svc.ListTasks()
	for _, t := range tasks {
		fmt.Println(t)
	}
	svc.AddTask("Task 3")

	percent, _ := svc.PercentDone()
	println("Percent done:", percent)

}
