package main

import (
	"fmt"

	"github.com/pragmaticbutton/todo/internal/domain/task"
	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/pragmaticbutton/todo/internal/storage/memory"
	"github.com/pragmaticbutton/todo/internal/utils"
)

func main() {
	// NOTE: The logic in main is only for testing purposes.
	storage := memory.New()
	listSvc := service.NewListService(storage, storage)
	err := listSvc.AddList(service.AddListInput{Description: "shopping"})
	if err != nil {
		panic(err)
	}
	// lists, _ := listSvc.ListLists()
	// for _, l := range lists {
	// 	fmt.Println(l)
	// }

	taskSvc := service.NewTaskService(storage, storage)
	_, err = taskSvc.AddTask(service.AddTaskInput{Description: "cookies", ListID: utils.Ptr(uint32(1)), Priority: utils.Ptr(task.PriorityMedium)})
	if err != nil {
		panic(err)
	}
	// tasks, _ := taskSvc.ListTasks()
	// for _, t := range tasks {
	// 	fmt.Println(t)
	// }
	_, err = taskSvc.AddTask(service.AddTaskInput{Description: "nesto", ListID: utils.Ptr(uint32(1)), Priority: utils.Ptr(task.PriorityHigh)})
	if err != nil {
		panic(err)
	}

	ls, _ := listSvc.ListTasks(uint32(1))
	for _, l := range ls {
		fmt.Println(l.ID, *l.ListID)
	}
}
