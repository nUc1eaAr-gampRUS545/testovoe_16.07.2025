package task

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testovoe_16.07.2025/configs"
	"testovoe_16.07.2025/pkg/utils"
)

type TaskServiceInterface interface {
	CreateTask() (*Task, error)
	AddFilesToTask(taskID string, urls []string) (*Task, error)
	GetTaskStatus(taskID string) (*Task, error)
}
type TaskService struct {
	config      *configs.Config
	tasks       map[string]*Task
	taskMutex   sync.Mutex
	activeTasks int
}

func NewTaskService(cfg *configs.Config) *TaskService {
	return &TaskService{
		config: cfg,
		tasks:  make(map[string]*Task),
	}
}

func (service *TaskService) CreateTask() (*Task, error) {
	service.taskMutex.Lock()
	defer service.taskMutex.Unlock()

	if service.activeTasks >= service.config.MaxConcurrentTasks {
		return nil, errors.New(ServerIsBusyError)
	}

	taskID := utils.GenerateID()
	task := &Task{
		ID:     taskID,
		Status: "created",
		Files:  make([]File, 0),
	}

	service.tasks[taskID] = task
	service.activeTasks++

	return task, nil
}

func (service *TaskService) AddFilesToTask(taskID string, urls []string) (*Task, error) {
	service.taskMutex.Lock()
	defer service.taskMutex.Unlock()

	task, exists := service.tasks[taskID]
	if !exists {
		return nil, errors.New(TaskNotFoundError)
	}

	if task.Status == "completed" {
		return nil, errors.New(TaskCompletedError)
	}

	if len(task.Files)+len(urls) > service.config.MaxFilesInTask {
		return nil, fmt.Errorf("exceeds maximum files per task (%d)", service.config.MaxFilesInTask)
	}

	task.Status = "processing"

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, url := range urls {
		if !utils.IsAllowedExtension(url, service.config.AllowedExtensions) {
			mutex.Lock()
			task.Errors = append(task.Errors, fmt.Sprintf("invalid file extension: %s", url))
			mutex.Unlock()
			continue
		}

		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			resp, err := http.Get(u)
			if err != nil {
				mutex.Lock()
				task.Errors = append(task.Errors, fmt.Sprintf("failed to download %s: %v", u, err))
				mutex.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				mutex.Lock()
				task.Errors = append(task.Errors, fmt.Sprintf("failed to download %s: status %d", u, resp.StatusCode))
				mutex.Unlock()
				return
			}

			fileName := filepath.Base(u)
			file := File{
				URL:  u,
				Name: fileName,
			}

			mutex.Lock()
			task.Files = append(task.Files, file)
			mutex.Unlock()
		}(url)
	}

	wg.Wait()

	if len(task.Files) >= service.config.MaxFilesInTask {
		go service.createArchive(task)
	}

	return task, nil

}

func (service *TaskService) GetTaskStatus(taskID string) (*Task, error) {
	service.taskMutex.Lock()
	defer service.taskMutex.Unlock()

	task, exists := service.tasks[taskID]
	if !exists {
		return nil, errors.New(TaskNotFoundError)
	}

	return task, nil
}

func (service *TaskService) createArchive(task *Task) {
	archiveName := fmt.Sprintf("archives/%s.zip", task.ID)
	err := os.MkdirAll("archives", os.ModePerm)
	if err != nil {
		task.Errors = append(task.Errors, fmt.Sprintf("failed to create archives directory: %v", err))
		return
	}

	file, err := os.Create(archiveName)
	if err != nil {
		task.Errors = append(task.Errors, fmt.Sprintf("failed to create archive file: %v", err))
		return
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, f := range task.Files {
		wg.Add(1)
		go func(file File) {
			defer wg.Done()

			resp, err := http.Get(file.URL)
			if err != nil {
				mutex.Lock()
				task.Errors = append(task.Errors, fmt.Sprintf("failed to download %s for archive: %v", file.URL, err))
				mutex.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				mutex.Lock()
				task.Errors = append(task.Errors, fmt.Sprintf("failed to download %s for archive: status %d", file.URL, resp.StatusCode))
				mutex.Unlock()
				return
			}

			zipFile, err := zipWriter.Create(file.Name)
			if err != nil {
				mutex.Lock()
				task.Errors = append(task.Errors, fmt.Sprintf("failed to create file in archive: %v", err))
				mutex.Unlock()
				return
			}

			_, err = io.Copy(zipFile, resp.Body)
			if err != nil {
				mutex.Lock()
				task.Errors = append(task.Errors, fmt.Sprintf("failed to write file to archive: %v", err))
				mutex.Unlock()
				return
			}
		}(f)
	}

	wg.Wait()

	service.taskMutex.Lock()
	defer service.taskMutex.Unlock()

	task.Status = "completed"
	task.Archive = fmt.Sprintf("/archives/%s", task.ID)
	service.activeTasks--
}
