package main

import (
	"fmt"
	"todo/internal/service"
	"todo/internal/storage/memory"
	"todo/internal/task"
)

func main() {
	svc := service.New(memory.New())
	svc.AddTask("Task 1", task.PriorityHigh)
	svc.AddTask("Task 2", task.PriorityLow)
	tasks, _ := svc.ListTasks()
	for _, t := range tasks {
		fmt.Println(t)
	}
	svc.CompleteTask(2)
	tasks, _ = svc.ListTasks()
	for _, t := range tasks {
		fmt.Println(t)
	}
	svc.AddTask("Task 3", task.PriorityHigh)

	percent, _ := svc.PercentDone()
	println("Percent done:", percent)

}
