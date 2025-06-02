package app

import (
	"errors"
	"log"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type TaskService interface {
	Save(t domain.Task) (domain.Task, error)
	Find(id uint64) (interface{}, error)

	//update
	FindAll(uId uint64, status *domain.TaskStatus, date *time.Time) ([]domain.Task, error)
	//update

	Update(t domain.Task) (domain.Task, error)
	Delete(id uint64) error

	//new
	UpdateStatus(taskID uint64, userID uint64, status domain.TaskStatus) (domain.Task, error)
	//new
}

type taskService struct {
	taskRepo database.TaskRepository
}

func NewTaskService(tr database.TaskRepository) TaskService {
	return taskService{
		taskRepo: tr,
	}
}

func (s taskService) Save(t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Save(t)
	if err != nil {
		log.Printf("taskService.Save(s.taskRepo.Save): %s", err)
		return domain.Task{}, err
	}

	return task, nil
}

func (s taskService) Find(id uint64) (interface{}, error) {
	task, err := s.taskRepo.Find(id)
	if err != nil {
		log.Printf("taskService.Find(s.taskRepo.Find): %s", err)
		return domain.Task{}, err
	}

	return task, nil
}

func (s taskService) FindAll(uId uint64, status *domain.TaskStatus, date *time.Time) ([]domain.Task, error) {
	tasks, err := s.taskRepo.FindAllTasks(uId, status, date)
	if err != nil {
		log.Printf("taskService.FindAll(s.taskRepo.FindAllTasks): %s", err)
		return nil, err
	}

	return tasks, nil
}

func (s taskService) Update(t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Update(t)
	if err != nil {
		log.Printf("taskService.Update(s.taskRepo.Update): %s", err)
		return domain.Task{}, err
	}

	return task, nil
}

func (s taskService) Delete(id uint64) error {
	err := s.taskRepo.Delete(id)
	if err != nil {
		log.Printf("taskService.Delete(s.taskRepo.Delete): %s", err)
		return err
	}

	return nil
}

func (s taskService) UpdateStatus(taskID uint64, userID uint64, status domain.TaskStatus) (domain.Task, error) {
	// знаходимо задачу
	task, err := s.taskRepo.Find(taskID)
	if err != nil {
		return domain.Task{}, err
	}

	// перевіряємо власника
	if task.UserId != userID {
		return domain.Task{}, errors.New("access denied")
	}

	task.Status = status

	return s.taskRepo.Update(task)
}
