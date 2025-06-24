package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vinibsi/todo-api/internal/entity"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Create(todo *entity.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) GetByID(id uint) (*entity.Todo, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Todo), args.Error(1)
}

func (m *MockTodoRepository) GetAll(limit, offset int) ([]entity.Todo, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]entity.Todo), args.Get(1).(int64), args.Error(2)
}

func (m *MockTodoRepository) Update(todo *entity.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTodoRepository) GetByCompleted(completed bool, limit, offset int) ([]entity.Todo, int64, error) {
	args := m.Called(completed, limit, offset)
	return args.Get(0).([]entity.Todo), args.Get(1).(int64), args.Error(2)
}
