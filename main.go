package main

import (
	"fmt"
	"todo/internal/storage"
)

func main() {

	storage := storage.New()
	storage.AddTask("prikolica")
	fmt.Println(storage.ListTasks())
}
