package ports

import (
	"context"
	"task_service/src/internal/core/dto"
)

type TaskRepository interface {
	CreateNewTask(ctx context.Context, taskData *dto.TaskDetails) (*dto.TaskDetails, error)
	UpdateExistingTask(ctx context.Context, taskID string, updatedData *dto.UpdatedTaskDetails) (*dto.TaskDetails, error)
	FetchTasks(ctx context.Context, userID, status string) ([]*dto.TaskDetails, error)
}
