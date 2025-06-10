package dto

import "time"

type TaskDetails struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Priority    string    `json:"priority,omitempty"`
	Status      string    `json:"status,omitempty"`
	AssignedTo  string    `json:"assigned_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewTaskDetails() *TaskDetails {
	now := time.Now()
	return &TaskDetails{
		Priority:  "medium",
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
