package usecase

import (
	"context"
	"task_service/src/internal/core/dto"
)

type Service interface {
	CreateTask(ctx context.Context, taskData *dto.TaskDetails) (*dto.TaskDetails, error)
	UpdatedTask(ctx context.Context, userId string, updatedData *dto.UpdatedTaskDetails) (*dto.TaskDetails, error)
	ListTasks(ctx context.Context, userID, status string) ([]*dto.TaskDetails, error)
}
