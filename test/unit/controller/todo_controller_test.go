package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vinibsi/todo-api/internal/controller"
	"github.com/vinibsi/todo-api/internal/dto"
	"github.com/vinibsi/todo-api/mocks"
)

type TodoControllerTestSuite struct {
	suite.Suite
	mockService    *mocks.MockTodoService
	todoController *controller.TodoController
	router         *gin.Engine
}

func (suite *TodoControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockService = new(mocks.MockTodoService)
	suite.todoController = controller.NewTodoController(suite.mockService)

	suite.router = gin.New()
	api := suite.router.Group("/api/v1")
	todos := api.Group("/todos")
	{
		todos.GET("", suite.todoController.GetAll)
		todos.GET("/:id", suite.todoController.GetByID)
		todos.POST("", suite.todoController.Create)
		todos.PUT("/:id", suite.todoController.Update)
		todos.DELETE("/:id", suite.todoController.Delete)
		todos.PATCH("/:id/complete", suite.todoController.Complete)
	}
}

func (suite *TodoControllerTestSuite) TestCreate_Success() {
	req := dto.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "high",
	}

	expectedResponse := &dto.TodoResponse{
		ID:          1,
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "high",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.mockService.On("Create", mock.AnythingOfType("*dto.CreateTodoRequest")).Return(expectedResponse, nil)

	jsonData, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, request)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Tarefa marcada como concluída", response.Message)
	suite.mockService.AssertExpectations(suite.T())
}

func TestTodoControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TodoControllerTestSuite))
}
