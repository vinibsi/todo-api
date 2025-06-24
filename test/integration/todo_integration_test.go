package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vinibsi/todo-api/internal/dto"
	"github.com/vinibsi/todo-api/internal/entity"
)

type TodoIntegrationTestSuite struct {
	suite.Suite
	helper *TestHelper
}

func (suite *TodoIntegrationTestSuite) SetupSuite() {
	helper, err := NewTestHelper()
	suite.Require().NoError(err)
	suite.helper = helper
}

func (suite *TodoIntegrationTestSuite) TearDownTest() {
	suite.helper.CleanDatabase()
}

func (suite *TodoIntegrationTestSuite) TestCreateTodo_Success() {
	req := dto.CreateTodoRequest{
		Title:       "Integration Test Todo",
		Description: "Testing todo creation",
		Priority:    "high",
	}

	httpReq, err := suite.helper.CreateTodoRequest("POST", "/api/v1/todos", req)
	suite.Require().NoError(err)

	recorder := suite.helper.ExecuteRequest(httpReq)

	assert.Equal(suite.T(), http.StatusCreated, recorder.Code)

	response, err := suite.helper.ParseSuccessResponse(recorder)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Tarefa criada com sucesso", response.Message)
	assert.NotNil(suite.T(), response.Data)

	// Verifica se o todo foi criado no banco
	var count int64
	suite.helper.DB.Model(&entity.Todo{}).Count(&count)
	assert.Equal(suite.T(), int64(1), count)
}

func (suite *TodoIntegrationTestSuite) TestCreateTodo_ValidationError() {
	req := dto.CreateTodoRequest{
		Title:       "", // Título vazio deve falhar na validação
		Description: "Testing validation",
	}

	httpReq, err := suite.helper.CreateTodoRequest("POST", "/api/v1/todos", req)
	suite.Require().NoError(err)

	recorder := suite.helper.ExecuteRequest(httpReq)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)

	response, err := suite.helper.ParseErrorResponse(recorder)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Dados inválidos", response.Error)
}

func (suite *TodoIntegrationTestSuite) TestGetTodoByID_Success() {
	// Primeiro cria um todo
	todo := &entity.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "medium",
	}
	suite.helper.Repository.Create(todo)

	// Busca o todo criado
	url := fmt.Sprintf("/api/v1/todos/%d", todo.ID)
	httpReq, err := suite.helper.CreateTodoRequest("GET", url, nil)
	suite.Require().NoError(err)

	recorder := suite.helper.ExecuteRequest(httpReq)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	response, err := suite.helper.ParseSuccessResponse(recorder)
	suite.Require().NoError(err)

	// Converte os dados para TodoResponse
	todoData, err := json.Marshal(response.Data)
	suite.Require().NoError(err)

	var todoResponse dto.TodoResponse
	err = json.Unmarshal(todoData, &todoResponse)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), todo.ID, todoResponse.ID)
	assert.Equal(suite.T(), "Test Todo", todoResponse.Title)
	assert.Equal(suite.T(), "medium", todoResponse.Priority)
}

func (suite *TodoIntegrationTestSuite) TestGetTodoByID_NotFound() {
	httpReq, err := suite.helper.CreateTodoRequest("GET", "/api/v1/todos/999", nil)
	suite.Require().NoError(err)

	recorder := suite.helper.ExecuteRequest(httpReq)

	assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)

	response, err := suite.helper.ParseErrorResponse(recorder)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Erro ao buscar tarefa", response.Error)
}

func (suite *TodoIntegrationTestSuite) TestGetAllTodos_Success() {
	// Cria múltiplos todos
	todos := []*entity.Todo{
		{Title: "Todo 1", Priority: "high"},
		{Title: "Todo 2", Priority: "medium"},
		{Title: "Todo 3", Priority: "low"},
	}

	for _, todo := range todos {
		suite.helper.Repository.Create(todo)
	}

	httpReq, err := suite.helper.CreateTodoRequest("GET", "/api/v1/todos?page=1&page_size=10", nil)
	suite.Require().NoError(err)

	recorder := suite.helper.ExecuteRequest(httpReq)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	response, err := suite.helper.ParseSuccessResponse(recorder)
	suite.Require().NoError(err)

	// Converte os dados para TodoListResponse
	listData, err := json.Marshal(response.Data)
	suite.Require().NoError(err)

	var listResponse dto.TodoListResponse
	err = json.Unmarshal(listData, &listResponse)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), int64(3), listResponse.Total)
	assert.Len(suite.T(), listResponse.Data, 3)
	assert.Equal(suite.T(), 1, listResponse.Page)
	assert.Equal(suite.T(), 10, listResponse.PageSize)
}

func (suite *TodoIntegrationTestSuite) TestUpdateTodo_Success() {
	// Cria um todo
	todo := &entity.Todo{
		Title:       "Original Title",
		Description: "Original Description",
		Priority:    "low",
		Completed:   false,
	}
	suite.helper.Repository.Create(todo)

	// Atualiza o todo
	newTitle := "Updated Title"
	completed := true
	updateReq := dto.UpdateTodoRequest{
		Title:     &newTitle,
		Completed: &completed,
	}

	url := fmt.Sprintf("/api/v1/todos/%d", todo.ID)
	httpReq, err := suite.helper.CreateTodoRequest("PUT", url, updateReq)
	suite.Require().NoError(err)

	recorder := suite.helper.ExecuteRequest(httpReq)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	response, err := suite.helper.ParseSuccessResponse(recorder)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Tarefa atualizada com sucesso", response.Message)

	// Verifica se foi atualizado no banco
	updatedTodo, err := suite.helper.Repository.GetByID(todo.ID)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Updated Title", updatedTodo.Title)
	assert.True(suite.T(), updatedTodo.Completed)
	assert.Equal(suite.T(), "Original Description", updatedTodo.Description) // Não deve ter mudado
}

func (suite *TodoIntegrationTestSuite) TestDeleteTodo_Success() {
	// Cria um todo
	todo := &entity.Todo{
		Title:    "To be deleted",
		Priority: "medium",
	}
	suite.helper.Repository.Create(todo)

	// Deleta o todo
	url := fmt.Sprintf("/api/v1/todos/%d", todo.ID)
	httpReq, err := suite.helper.CreateTodoRequest("DELETE", url, nil)
	suite.Require().NoError(err)

	recorder := suite.helper.ExecuteRequest(httpReq)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	response, err := suite.helper.ParseSuccessResponse(recorder)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Tarefa deletada com sucesso", response.Message)

	// Verifica se foi deletado do banco (soft delete)
	_, err = suite.helper.Repository.GetByID(todo.ID)
	assert.Error(suite.T(), err) // Deve retornar erro pois foi deletado
}

func (suite *TodoIntegrationTestSuite) TestCompleteTodo_Success() {
	// Cria um todo não concluído
	todo := &entity.Todo{
		Title:     "To be completed",
		Priority:  "high",
		Completed: false,
	}
	suite.helper.Repository.Create(todo)

	// Marca como concluído
	url := fmt.Sprintf("/api/v1/todos/%d/complete", todo.ID)
	httpReq, err := suite.helper.CreateTodoRequest("PATCH", url, nil)
	suite.Require().NoError(err)

	recorder := suite.helper.ExecuteRequest(httpReq)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	response, err := suite.helper.ParseSuccessResponse(recorder)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Tarefa marcada como concluída", response.Message)

	// Verifica se foi marcado como concluído no banco
	completedTodo, err := suite.helper.Repository.GetByID(todo.ID)
	suite.Require().NoError(err)

	assert.True(suite.T(), completedTodo.Completed)
}

func (suite *TodoIntegrationTestSuite) TestFullWorkflow() {
	// Teste de workflow completo: criar, buscar, atualizar, marcar como concluída e deletar

	// 1. Criar todo
	createReq := dto.CreateTodoRequest{
		Title:       "Workflow Todo",
		Description: "Testing full workflow",
		Priority:    "high",
	}

	httpReq, _ := suite.helper.CreateTodoRequest("POST", "/api/v1/todos", createReq)
	recorder := suite.helper.ExecuteRequest(httpReq)
	assert.Equal(suite.T(), http.StatusCreated, recorder.Code)

	// Extrai o ID do todo criado
	response, _ := suite.helper.ParseSuccessResponse(recorder)
	todoData, _ := json.Marshal(response.Data)
	var createdTodo dto.TodoResponse
	json.Unmarshal(todoData, &createdTodo)
	todoID := createdTodo.ID

	// 2. Buscar todo
	url := fmt.Sprintf("/api/v1/todos/%d", todoID)
	httpReq, _ = suite.helper.CreateTodoRequest("GET", url, nil)
	recorder = suite.helper.ExecuteRequest(httpReq)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	// 3. Atualizar todo
	newTitle := "Updated Workflow Todo"
	updateReq := dto.UpdateTodoRequest{Title: &newTitle}
	httpReq, _ = suite.helper.CreateTodoRequest("PUT", url, updateReq)
	recorder = suite.helper.ExecuteRequest(httpReq)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	// 4. Marcar como concluída
	completeURL := fmt.Sprintf("/api/v1/todos/%d/complete", todoID)
	httpReq, _ = suite.helper.CreateTodoRequest("PATCH", completeURL, nil)
	recorder = suite.helper.ExecuteRequest(httpReq)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	// 5. Deletar todo
	httpReq, _ = suite.helper.CreateTodoRequest("DELETE", url, nil)
	recorder = suite.helper.ExecuteRequest(httpReq)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	// 6. Verificar que foi deletado
	httpReq, _ = suite.helper.CreateTodoRequest("GET", url, nil)
	recorder = suite.helper.ExecuteRequest(httpReq)
	assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)
}

func TestTodoIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(TodoIntegrationTestSuite))
}
