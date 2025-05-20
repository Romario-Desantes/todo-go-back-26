package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskDto struct {
	Id          uint64            `json:"id"`
	UserId      uint64            `json:"userId"`
	Title       string            `json:"title"`
	Description *string           `json:"description,omitempty"`
	Date        *time.Time        `json:"date,omitempty"`
	Status      domain.TaskStatus `json:"status"`
}

func (d TaskDto) DomainToDto(t domain.Task) TaskDto {
	return TaskDto{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Date:        t.Date,
		Status:      t.Status,
	}
}

func (d TaskDto) DomainToDtoCollection(ts []domain.Task) []TaskDto {
	tasksDto := make([]TaskDto, len(ts))
	for i, t := range ts {
		tasksDto[i] = d.DomainToDto(t)
	}

	return tasksDto
}
