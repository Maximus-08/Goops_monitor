package runner

import (
	"time"
)

type Task struct {
	ID        string
	Script    string
	CreatedAt time.Time
	Status    string
}

func NewTask(id, script string) *Task {
	return &Task{
		ID:        id,
		Script:    script,
		CreatedAt: time.Now(),
		Status:    "pending",
	}
}

func (t *Task) MarkRunning() {
	t.Status = "running"
}

func (t *Task) MarkCompleted() {
	t.Status = "completed"
}

func (t *Task) MarkFailed() {
	t.Status = "failed"
}
