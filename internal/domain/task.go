package domain

import "time"

type Task struct {
	Id          uint64
	UserId      uint64
	Title       string
	Description *string
	Date        *time.Time
	Status      string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}
