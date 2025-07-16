package main

import (
	"log"
	"net/http"
	"testovoe_16.07.2025/configs"
	"testovoe_16.07.2025/internal/task"
	"testovoe_16.07.2025/middleware"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()

	//Services
	taskService := task.NewTaskService(config)

	//Handlers
	task.NewTaskHandler(
		router,
		task.TaskHandlerDeps{
			TaskService: taskService,
		},
	)

	stackMiddleware := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    config.Port,
		Handler: stackMiddleware(router),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
