package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

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

		// Фільтр по статусу (опціональний)
		var status *domain.TaskStatus
		statusStr := r.URL.Query().Get("status")
		if statusStr != "" {
			st := domain.TaskStatus(statusStr)
			status = &st
		}

		// Фільтр по даті (опціональний)
		var date *time.Time
		dateStr := r.URL.Query().Get("date")
		if dateStr != "" {
			parsedDate, err := time.Parse("2006-01-02", dateStr) // формат: YYYY-MM-DD
			if err != nil {
				BadRequest(w, errors.New("invalid date format (use YYYY-MM-DD)"))
				return
			}
			date = &parsedDate
		}

		// Виклик сервісу з фільтрами
		tasks, err := c.taskService.FindAll(user.Id, status, date)
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

func (c TaskController) UpdateStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отримуєм ID задачі
		task := r.Context().Value(TaskKey).(domain.Task)
		user := r.Context().Value(UserKey).(domain.User)

		// Парсимо статус із JSON запиту
		var body struct {
			Status domain.TaskStatus `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			BadRequest(w, err)
			return
		}

		// Викликаємо сервіс з частковим оновленням ресурсу
		updatedTask, err := c.taskService.UpdateStatus(task.Id, user.Id, body.Status)
		if err != nil {
			if err.Error() == "access denied" {
				Forbidden(w, err)
				return
			}
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(updatedTask)
		Success(w, taskDto)
	}
}
