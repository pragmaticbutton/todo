package main

import (
	"fmt"
	"todo/internal/domain/list"
	"todo/internal/domain/task"
	"todo/internal/service"
	"todo/internal/storage/memory"
	"todo/internal/utils"
)

func main() {
	storage := memory.New()
	svc := service.NewTaskService(storage)
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

	svc.UpdateTask(3, service.UpdateTaskInput{Description: utils.Ptr("Novi description")})
	t, _ := svc.GetTask(3)
	fmt.Println(t)

	fmt.Println("-----------------------")
	storage.AddList(&list.List{ID: 1, Description: "List 1"})
	ls, _ := storage.ListLists()
	for _, l := range ls {
		fmt.Println(l)
	}
}
