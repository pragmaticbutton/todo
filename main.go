package main

import (
	"fmt"
	"todo/internal/storage"
	"todo/internal/task"
)

func main() {

	storage := storage.New()
	storage.AddTask(&task.Task{
		ID:          storage.NextID(),
		Description: "prvi zadatak"})
	t, _ := storage.GetTask(1)
	fmt.Println(*t)
	t.Description = "novi opis"
	t.Done = true
	storage.UpdateTask(t)
	t, _ = storage.GetTask(1)
	fmt.Println(*t)
}
