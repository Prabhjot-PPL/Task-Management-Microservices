package persistance

import (
	"context"
	"fmt"
	"task_service/src/internal/adaptors/ports"
	"task_service/src/internal/core/dto"
	"task_service/src/pkg/logger"
)

type TaskRepo struct {
	db *Database
}

func NewTaskRepo(d *Database) ports.TaskRepository {
	return &TaskRepo{db: d}
}

func (t *TaskRepo) CreateNewTask(ctx context.Context, taskData *dto.TaskDetails) (*dto.TaskDetails, error) {

	task := dto.NewTaskDetails()
	var taskId int

	query1 := `INSERT INTO task_details (title, description, priority, status, assigned_to)
				VALUES($1, $2, $3, $4, $5)
				RETURNING id`
	err := t.db.db.QueryRowContext(ctx, query1, taskData.Title, taskData.Description, taskData.Priority, taskData.Status, taskData.AssignedTo).Scan(&taskId)
	if err != nil {
		fmt.Printf("\n")
		logger.Log.Error("Error creating task : ", err)
		return task, err
	}

	query2 := `SELECT id, title, description, priority, status, created_at, updated_at, assigned_to FROM task_details WHERE id=$1`
	row := t.db.db.QueryRowContext(ctx, query2, taskId)

	err = row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.CreatedAt, &task.UpdatedAt, &task.AssignedTo)
	if err != nil {
		fmt.Printf("\n")
		logger.Log.Error("Error creating task : ", err)
		return task, err
	}

	return task, nil
}

func (t *TaskRepo) UpdateExistingTask(ctx context.Context, taskID string, updatedData *dto.UpdatedTaskDetails) (*dto.TaskDetails, error) {
	var existing dto.TaskDetails

	// 1. Fetch existing task
	querySelect := `SELECT title, description, priority, status, assigned_to FROM task_details WHERE id = $1`
	err := t.db.db.QueryRowContext(ctx, querySelect, taskID).Scan(
		&existing.Title,
		&existing.Description,
		&existing.Priority,
		&existing.Status,
		&existing.AssignedTo,
	)
	if err != nil {
		logger.Log.Error("Error fetching existing task: ", err)
		return nil, err
	}

	// 2. Map fields: prefer updatedData if provided (non-zero/empty)
	if updatedData.Title != "" {
		existing.Title = updatedData.Title
	}
	if updatedData.Description != "" {
		existing.Description = updatedData.Description
	}
	if updatedData.Priority != "" {
		existing.Priority = updatedData.Priority
	}
	if updatedData.Status != "" {
		existing.Status = updatedData.Status
	}
	if updatedData.AssignedTo != "" {
		existing.AssignedTo = updatedData.AssignedTo
	}

	// 3. Update in DB
	queryUpdate := `UPDATE task_details SET title=$1, description=$2, priority=$3, status=$4, assigned_to=$5, updated_at=NOW() WHERE id=$6`
	_, err = t.db.db.ExecContext(ctx, queryUpdate,
		existing.Title,
		existing.Description,
		existing.Priority,
		existing.Status,
		existing.AssignedTo,
		taskID,
	)
	if err != nil {
		logger.Log.Error("Error updating task: ", err)
		return nil, err
	}

	return &existing, nil
}

func (t *TaskRepo) FetchTasks(ctx context.Context, userID, status string) ([]*dto.TaskDetails, error) {
	query := `SELECT id, title, description, priority, status, created_at, updated_at, assigned_to FROM task_details WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if userID != "" {
		query += fmt.Sprintf(" AND assigned_to = $%d", argPos)
		args = append(args, userID)
		argPos++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, status)
		argPos++
	}

	rows, err := t.db.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Error fetching tasks: ", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []*dto.TaskDetails
	for rows.Next() {
		task := dto.NewTaskDetails()
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.AssignedTo,
		)
		if err != nil {
			logger.Log.Error("Error scanning task row: ", err)
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
