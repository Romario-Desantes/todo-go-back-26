package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.RegisterRequest{}, domain.User{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, errors.New("invalid request body"))
			return
		}

		user, token, err := c.authService.Register(user)
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}

		var authDto resources.AuthDto
		Success(w, authDto.DomainToDto(token, user))
	}
}
