package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vinibsi/todo-api/internal/dto"
)

type MockTodoService struct {
	mock.Mock
}

func (m *MockTodoService) Create(req *dto.CreateTodoRequest) (*dto.TodoResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.TodoResponse), args.Error(1)
}

func (m *MockTodoService) GetByID(id uint) (*dto.TodoResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.TodoResponse), args.Error(1)
}

func (m *MockTodoService) GetAll(page, pageSize int) (*dto.TodoListResponse, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).(*dto.TodoListResponse), args.Error(1)
}

func (m *MockTodoService) Update(id uint, req *dto.UpdateTodoRequest) (*dto.TodoResponse, error) {
	args := m.Called(id, req)
	return args.Get(0).(*dto.TodoResponse), args.Error(1)
}

func (m *MockTodoService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTodoService) Complete(id uint) (*dto.TodoResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.TodoResponse), args.Error(1)
}
