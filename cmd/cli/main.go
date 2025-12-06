package main

import (
	"log"

	"github.com/pragmaticbutton/todo/internal/cli"
	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/pragmaticbutton/todo/internal/storage/memory"
)

func main() {
	storage := memory.New()
	taskService := service.NewTaskService(storage, storage)
	listService := service.NewListService(storage, storage)

	rootCmd := cli.NewRootCmd(taskService, listService)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
