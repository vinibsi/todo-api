package service

import (
	"errors"
	"math"

	"github.com/vinibsi/todo-api/internal/dto"
	"github.com/vinibsi/todo-api/internal/entity"
	"github.com/vinibsi/todo-api/internal/repository"
	"gorm.io/gorm"
)

type TodoService interface {
	Create(req *dto.CreateTodoRequest) (*dto.TodoResponse, error)
	GetByID(id uint) (*dto.TodoResponse, error)
	GetAll(page, pageSize int) (*dto.TodoListResponse, error)
	Update(id uint, req *dto.UpdateTodoRequest) (*dto.TodoResponse, error)
	Delete(id uint) error
	Complete(id uint) (*dto.TodoResponse, error)
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) Create(req *dto.CreateTodoRequest) (*dto.TodoResponse, error) {
	todo := &entity.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Completed:   false,
	}

	if todo.Priority == "" {
		todo.Priority = "medium"
	}

	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}

	return s.entityToDTO(todo), nil
}

func (s *todoService) GetByID(id uint) (*dto.TodoResponse, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}
	return s.entityToDTO(todo), nil
}

func (s *todoService) Complete(id uint) (*dto.TodoResponse, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	todo.Completed = true
	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	return s.entityToDTO(todo), nil
}

func (s *todoService) GetAll(page, pageSize int) (*dto.TodoListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	todos, total, err := s.repo.GetAll(pageSize, offset)
	if err != nil {
		return nil, err
	}

	todoResponses := make([]dto.TodoResponse, len(todos))
	for i, todo := range todos {
		todoResponses[i] = *s.entityToDTO(&todo)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.TodoListResponse{
		Data:       todoResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *todoService) Update(id uint, req *dto.UpdateTodoRequest) (*dto.TodoResponse, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	// Atualiza somente campos fornecidos
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	return s.entityToDTO(todo), nil
}

func (s *todoService) Delete(id uint) error {
	// Verifica se a tarefa existe
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("todo not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *todoService) entityToDTO(todo *entity.Todo) *dto.TodoResponse {
	return &dto.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		Priority:    todo.Priority,
		DueDate:     todo.DueDate,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}
