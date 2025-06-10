package dto

type TaskNotification struct {
	Event      string `json:"event"` // "assigned" or "updated"
	TaskID     string `json:"id"`
	AssignedTo string `json:"assigned_to"`
	Message    string `json:"message"`
}
