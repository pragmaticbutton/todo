package main

import (
	"fmt"
	"todo/internal/storage"
)

func main() {

	storage := storage.New()
	storage.AddTask("prikolica")
	t, _ := storage.GetTask(1)
	fmt.Println(*t)
	t.Description = "novi opis"
	t.Done = true
	storage.UpdateTask(t)
	t, _ = storage.GetTask(1)
	fmt.Println(*t)
}
