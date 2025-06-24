package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vinibsi/todo-api/internal/dto"
	"github.com/vinibsi/todo-api/internal/entity"
	"github.com/vinibsi/todo-api/internal/service"
	"github.com/vinibsi/todo-api/mocks"
	"gorm.io/gorm"
)

type TodoServiceTestSuite struct {
	suite.Suite
	mockRepo    *mocks.MockTodoRepository
	todoService service.TodoService
}

func (suite *TodoServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockTodoRepository)
	suite.todoService = service.NewTodoService(suite.mockRepo)
}

func (suite *TodoServiceTestSuite) TestCreate_Success() {
	req := &dto.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "high",
	}

	expectedTodo := &entity.Todo{
		ID:          1,
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "high",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.mockRepo.On("Create", mock.AnythingOfType("*entity.Todo")).Return(nil).Run(func(args mock.Arguments) {
		todo := args.Get(0).(*entity.Todo)
		todo.ID = 1
		todo.CreatedAt = expectedTodo.CreatedAt
		todo.UpdatedAt = expectedTodo.UpdatedAt
	})

	result, err := suite.todoService.Create(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), uint(1), result.ID)
	assert.Equal(suite.T(), "Test Todo", result.Title)
	assert.Equal(suite.T(), "high", result.Priority)
	assert.False(suite.T(), result.Completed)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TodoServiceTestSuite) TestCreate_DefaultPriority() {
	req := &dto.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		// Priority não informado
	}

	suite.mockRepo.On("Create", mock.AnythingOfType("*entity.Todo")).Return(nil).Run(func(args mock.Arguments) {
		todo := args.Get(0).(*entity.Todo)
		todo.ID = 1
		assert.Equal(suite.T(), "medium", todo.Priority) // Verifica prioridade padrão
	})

	result, err := suite.todoService.Create(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "medium", result.Priority)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TodoServiceTestSuite) TestGetByID_Success() {
	todo := &entity.Todo{
		ID:          1,
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "high",
		Completed:   false,
	}

	suite.mockRepo.On("GetByID", uint(1)).Return(todo, nil)

	result, err := suite.todoService.GetByID(1)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), uint(1), result.ID)
	assert.Equal(suite.T(), "Test Todo", result.Title)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TodoServiceTestSuite) TestGetByID_NotFound() {
	suite.mockRepo.On("GetByID", uint(999)).Return((*entity.Todo)(nil), gorm.ErrRecordNotFound)

	result, err := suite.todoService.GetByID(999)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "todo not found", err.Error())
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TodoServiceTestSuite) TestGetAll_Success() {
	todos := []entity.Todo{
		{ID: 1, Title: "Todo 1", Priority: "high"},
		{ID: 2, Title: "Todo 2", Priority: "medium"},
	}

	suite.mockRepo.On("GetAll", 10, 0).Return(todos, int64(2), nil)

	result, err := suite.todoService.GetAll(1, 10)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Data, 2)
	assert.Equal(suite.T(), int64(2), result.Total)
	assert.Equal(suite.T(), 1, result.Page)
	assert.Equal(suite.T(), 10, result.PageSize)
	assert.Equal(suite.T(), 1, result.TotalPages)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TodoServiceTestSuite) TestUpdate_Success() {
	existingTodo := &entity.Todo{
		ID:          1,
		Title:       "Original Title",
		Description: "Original Description",
		Priority:    "low",
		Completed:   false,
	}

	newTitle := "Updated Title"
	completed := true
	req := &dto.UpdateTodoRequest{
		Title:     &newTitle,
		Completed: &completed,
	}

	suite.mockRepo.On("GetByID", uint(1)).Return(existingTodo, nil)
	suite.mockRepo.On("Update", mock.AnythingOfType("*entity.Todo")).Return(nil).Run(func(args mock.Arguments) {
		todo := args.Get(0).(*entity.Todo)
		assert.Equal(suite.T(), "Updated Title", todo.Title)
		assert.True(suite.T(), todo.Completed)
		assert.Equal(suite.T(), "Original Description", todo.Description) // Não deve mudar
	})

	result, err := suite.todoService.Update(1, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TodoServiceTestSuite) TestDelete_Success() {
	todo := &entity.Todo{ID: 1, Title: "To be deleted"}

	suite.mockRepo.On("GetByID", uint(1)).Return(todo, nil)
	suite.mockRepo.On("Delete", uint(1)).Return(nil)

	err := suite.todoService.Delete(1)

	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TodoServiceTestSuite) TestComplete_Success() {
	todo := &entity.Todo{
		ID:        1,
		Title:     "Test Todo",
		Completed: false,
	}

	suite.mockRepo.On("GetByID", uint(1)).Return(todo, nil)
	suite.mockRepo.On("Update", mock.AnythingOfType("*entity.Todo")).Return(nil).Run(func(args mock.Arguments) {
		updatedTodo := args.Get(0).(*entity.Todo)
		assert.True(suite.T(), updatedTodo.Completed)
	})

	result, err := suite.todoService.Complete(1)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.True(suite.T(), result.Completed)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestTodoServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TodoServiceTestSuite))
}
