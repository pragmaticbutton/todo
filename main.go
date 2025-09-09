package main

import (
	"fmt"
	"todo/internal/storage"
)

func main() {

	storage := storage.New()
	storage.AddTask("prikolica")
	tasks, _ := storage.ListTasks()
	fmt.Println(tasks)
}
