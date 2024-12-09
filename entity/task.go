package entity

const (
	StatusPending   Status = "Pending"
	StatusCompleted Status = "Completed"
)

type Status string
type Task struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}
