package dto

type UpdatedTaskDetails struct {
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    string `json:"priority,omitempty"`
	Status      string `json:"status,omitempty"`
	AssignedTo  string `json:"assigned_to,omitempty"`
}
