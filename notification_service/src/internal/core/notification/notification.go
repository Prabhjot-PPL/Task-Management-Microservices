package notification

type TaskNotification struct {
	Event      string `json:"event"`
	TaskID     string `json:"task_id"`
	AssignedTo string `json:"assigned_to"`
	Message    string `json:"message"`
}
