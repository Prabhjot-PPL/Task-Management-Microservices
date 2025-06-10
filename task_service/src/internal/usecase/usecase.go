package usecase

import (
	"context"
	"task_service/src/internal/adaptors/ports"
	"task_service/src/internal/core/dto"
)

type TaskService struct {
	taskRepo ports.TaskRepository
}

func NewTaskService(taskRepo ports.TaskRepository) Service {
	return &TaskService{taskRepo: taskRepo}
}

func (t *TaskService) CreateTask(ctx context.Context, taskData *dto.TaskDetails) (*dto.TaskDetails, error) {
	return t.taskRepo.CreateNewTask(ctx, taskData)
}

func (t *TaskService) UpdatedTask(ctx context.Context, taskID string, updatedData *dto.UpdatedTaskDetails) (*dto.TaskDetails, error) {
	return t.taskRepo.UpdateExistingTask(ctx, taskID, updatedData)
}

func (t *TaskService) ListTasks(ctx context.Context, userID, status string) ([]*dto.TaskDetails, error) {
	return t.taskRepo.FetchTasks(ctx, userID, status)
}
