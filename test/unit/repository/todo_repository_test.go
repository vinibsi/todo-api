package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vinibsi/todo-api/internal/entity"
	"github.com/vinibsi/todo-api/internal/repository"
	"github.com/vinibsi/todo-api/pkg/database"
	"gorm.io/gorm"
)

type TodoRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo repository.TodoRepository
}

func (suite *TodoRepositoryTestSuite) SetupSuite() {
	db, err := database.ConnectTest()
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = repository.NewTodoRepository(db)
}

func (suite *TodoRepositoryTestSuite) TearDownTest() {
	// Limpa a tabela após cada teste
	suite.db.Exec("DELETE FROM todos")
}

func (suite *TodoRepositoryTestSuite) TestCreate() {
	todo := &entity.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "high",
		Completed:   false,
	}

	err := suite.repo.Create(todo)

	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), todo.ID)
	assert.NotZero(suite.T(), todo.CreatedAt)
}

func (suite *TodoRepositoryTestSuite) TestGetByID() {
	// Cria um todo
	todo := &entity.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "medium",
	}
	suite.repo.Create(todo)

	// Busca o todo
	found, err := suite.repo.GetByID(todo.ID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), todo.Title, found.Title)
	assert.Equal(suite.T(), todo.Description, found.Description)
	assert.Equal(suite.T(), todo.Priority, found.Priority)
}

func (suite *TodoRepositoryTestSuite) TestGetAll() {
	// Cria múltiplos todos
	todos := []*entity.Todo{
		{Title: "Todo 1", Priority: "high"},
		{Title: "Todo 2", Priority: "medium"},
		{Title: "Todo 3", Priority: "low"},
	}

	for _, todo := range todos {
		suite.repo.Create(todo)
	}

	// Busca todos
	result, total, err := suite.repo.GetAll(10, 0)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), total)
	assert.Len(suite.T(), result, 3)
}

func (suite *TodoRepositoryTestSuite) TestUpdate() {
	// Cria um todo
	todo := &entity.Todo{
		Title:     "Original Title",
		Priority:  "low",
		Completed: false,
	}
	suite.repo.Create(todo)

	// Atualiza o todo
	todo.Title = "Updated Title"
	todo.Completed = true
	err := suite.repo.Update(todo)

	assert.NoError(suite.T(), err)

	// Verifica a atualização
	updated, _ := suite.repo.GetByID(todo.ID)
	assert.Equal(suite.T(), "Updated Title", updated.Title)
	assert.True(suite.T(), updated.Completed)
}

func (suite *TodoRepositoryTestSuite) TestDelete() {
	// Cria um todo
	todo := &entity.Todo{
		Title:    "To be deleted",
		Priority: "medium",
	}
	suite.repo.Create(todo)

	// Deleta o todo
	err := suite.repo.Delete(todo.ID)
	assert.NoError(suite.T(), err)

	// Verifica se foi deletado (soft delete)
	_, err = suite.repo.GetByID(todo.ID)
	assert.Error(suite.T(), err)
}

func TestTodoRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TodoRepositoryTestSuite))
}
