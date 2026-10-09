package todo

import "time"

type Task struct {
	Title       string
	Description string
	Completed   bool

	CreatedAt   time.Time
	CompletedAt *time.Time
}

func NewTask(title string, description string) Task {
	return Task{
		Title:       title,
		Description: description,

		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
}

func (t *Task) Complete() {
	t.Completed = true
	now := time.Now()
	t.CompletedAt = &now
}

func (t *Task) Uncomplete() {
	t.Completed = false
	t.CompletedAt = nil
}
