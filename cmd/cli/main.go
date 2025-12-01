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

	rootCmd := cli.NewRootCmd(taskService)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
