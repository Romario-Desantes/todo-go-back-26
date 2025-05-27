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
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController: %s", err)
			BadRequest(w, err)
			return
		}

		task.Status = domain.TaskNew
		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id

		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := r.Context().Value(TaskKey).(domain.Task)
		user := r.Context().Value(UserKey).(domain.User)

		if task.UserId != user.Id {
			err := errors.New("access denied")
			Forbidden(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		tasks, err := c.taskService.FindAll(user.Id)
		if err != nil {
			log.Printf("TaskController.FindAll(c.taskService.FindAll): %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		tasksDto := taskDto.DomainToDtoCollection(tasks)
		Success(w, tasksDto)
	}
}

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		taskExists := r.Context().Value(TaskKey).(domain.Task)
		if taskExists.UserId != user.Id {
			err = errors.New("access denied")
			Forbidden(w, err)
			return
		}

		taskExists.Title = task.Title
		taskExists.Description = task.Description
		taskExists.Date = task.Date

		task, err = c.taskService.Update(taskExists)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := r.Context().Value(TaskKey).(domain.Task)
		user := r.Context().Value(UserKey).(domain.User)

		if task.UserId != user.Id {
			err := errors.New("access denied")
			Forbidden(w, err)
			return
		}

		err := c.taskService.Delete(task.Id)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}
