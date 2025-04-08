package requests

import "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

type TaskRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	Date        *int64  `json:"date"`
}

func (r TaskRequest) ToDomainModel() (domain.Task, error) {
	date := 
}
