package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"task_service/src/internal/adaptors/persistance"
	"task_service/src/internal/core/dto"
	"task_service/src/internal/usecase"
	errorhandling "task_service/src/pkg/error_handling"
	pkgresponse "task_service/src/pkg/response"
	"time"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService usecase.Service
	publisher   *persistance.RedisPublisher
}

func NewTaskHandler(taskService usecase.Service, pub *persistance.RedisPublisher) *TaskHandler {
	return &TaskHandler{taskService: taskService, publisher: pub}
}

func (t *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	taskData := dto.NewTaskDetails()

	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		errorhandling.HandlerError(w, "failed to create task", http.StatusBadRequest, err)
		return
	}

	result, err := t.taskService.CreateTask(ctx, taskData)
	if err != nil {
		errorhandling.HandlerError(w, "failed to create task", http.StatusInternalServerError, err)
		return
	}

	if taskData.AssignedTo != "" {
		_ = t.publisher.Publish(ctx, "task_notifications", dto.TaskNotification{
			Event:      "assigned",
			TaskID:     taskData.ID,
			AssignedTo: *&taskData.AssignedTo,
			Message:    fmt.Sprintf("Task %s assigned to user %d", taskData.ID, *&taskData.AssignedTo),
		})
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Data:    result,
		Message: "User Registered Successfully ",
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (t *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	taskID := chi.URLParam(r, "id")
	var updateData dto.UpdatedTaskDetails

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		errorhandling.HandlerError(w, "Invalid request body", http.StatusInternalServerError, err)
		return
	}

	result, err := t.taskService.UpdatedTask(ctx, taskID, &updateData)
	if err != nil {
		errorhandling.HandlerError(w, "failed to update task", http.StatusInternalServerError, err)
		return
	}

	if result.AssignedTo != "" {
		_ = t.publisher.Publish(ctx, "task_notifications", dto.TaskNotification{
			Event:      "assigned",
			TaskID:     result.ID,
			AssignedTo: *&result.AssignedTo,
			Message:    fmt.Sprintf("Task %d assigned to user %d", result.ID, *&result.AssignedTo),
		})
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Data:    result,
		Message: "Task Updated Successfully ",
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (t *TaskHandler) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userID := r.URL.Query().Get("user_id")
	status := r.URL.Query().Get("status")

	tasks, err := t.taskService.ListTasks(ctx, userID, status)
	if err != nil {
		errorhandling.HandlerError(w, "failed to fetch tasks", http.StatusInternalServerError, err)
		return
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Data:    tasks,
		Message: "Tasks fetched successfully",
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}
