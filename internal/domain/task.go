package domain

import "time"

type Task struct {
	Id          uint64
	UserId      uint64
	Title       string
	Description *string
	Date        *time.Time
	Status      TaskStatus
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type TaskStatus string

const (
	TaskNew        TaskStatus = "NEW"
	TaskInProgress TaskStatus = "IN_PROGRESS"
	TaskComplete   TaskStatus = "COMPLETE"
)
