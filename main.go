package main

import (
	"fmt"
	"todo/internal/domain/task"
	"todo/internal/service"
	"todo/internal/storage/memory"
	"todo/internal/utils"
)

func main() {
	storage := memory.New()
	listSvc := service.NewListService(storage)
	listSvc.AddList(service.AddListInput{Description: "shopping"})
	lists, _ := listSvc.ListLists()
	for _, l := range lists {
		fmt.Println(l)
	}

	taskSvc := service.NewTaskService(storage)
	_, err := taskSvc.AddTask(service.AddTaskInput{Description: "cookies", ListID: utils.Ptr(uint32(3)), Priority: task.PriorityMedium})
	if err != nil {
		panic(err)
	}
	tasks, _ := taskSvc.ListTasks()
	for _, t := range tasks {
		fmt.Println(t)
	}
}
