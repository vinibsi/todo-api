package service

import (
	"github.com/vinibsi/todo-api/internal/dto"
	"github.com/vinibsi/todo-api/internal/entity"
	"github.com/vinibsi/todo-api/internal/repository"
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
