package task

import (
	"net/http"
	"testovoe_16.07.2025/pkg/request"
	"testovoe_16.07.2025/pkg/response"
)

type TaskHandlerDeps struct {
	TaskService TaskServiceInterface
}
type TaskHandler struct {
	TaskService TaskServiceInterface
}

type TaskHandlerInterface interface {
	CreateTask() http.HandlerFunc
	AddFilesToTask() http.HandlerFunc
	GetTaskStatus() http.HandlerFunc
}

func NewTaskHandler(router *http.ServeMux, deps TaskHandlerDeps) {

	handler := &TaskHandler{
		TaskService: deps.TaskService,
	}
	router.HandleFunc("POST /task", handler.CreateTask())
	router.HandleFunc("GET /task/{id}", handler.GetTaskStatus())
	router.HandleFunc("PUT /task/{id}", handler.AddFilesToTask())

}

func (handler *TaskHandler) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		task, err := handler.TaskService.CreateTask()
		if err != nil {
			response.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, task, http.StatusOK)
	}
}

func (handler *TaskHandler) GetTaskStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		taskId := req.PathValue("id")
		if taskId == "" {
			response.Json(w, "", http.StatusBadRequest)
			return
		}
		task, err := handler.TaskService.GetTaskStatus(taskId)
		if err != nil {
			response.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, task, http.StatusOK)
	}
}
func (handler *TaskHandler) AddFilesToTask() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		taskId := req.PathValue("id")
		if taskId == "" {
			response.Json(w, "", http.StatusBadRequest)
			return
		}
		body, err := request.HandlerBody[TaskAddFilesRequest](&w, req)
		if err != nil {
			response.Json(w, "", http.StatusBadRequest)
			return
		}
		task, err := handler.TaskService.AddFilesToTask(taskId, body.URLs)
		if err != nil {
			response.Json(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Json(w, task, http.StatusOK)
	}
}
