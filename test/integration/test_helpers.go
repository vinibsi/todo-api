package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/vinibsi/todo-api/internal/controller"
	"github.com/vinibsi/todo-api/internal/dto"
	"github.com/vinibsi/todo-api/internal/repository"
	"github.com/vinibsi/todo-api/internal/service"
	"github.com/vinibsi/todo-api/pkg/database"
	"gorm.io/gorm"
)

type TestHelper struct {
	DB         *gorm.DB
	Router     *gin.Engine
	Controller *controller.TodoController
	Service    service.TodoService
	Repository repository.TodoRepository
}

func NewTestHelper() (*TestHelper, error) {
	// Conecta ao banco de teste (SQLite em memória)
	db, err := database.ConnectTest()
	if err != nil {
		return nil, err
	}

	// Inicializa as camadas
	repo := repository.NewTodoRepository(db)
	svc := service.NewTodoService(repo)
	ctrl := controller.NewTodoController(svc)

	// Configura o router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api/v1")
	todos := api.Group("/todos")
	{
		todos.GET("", ctrl.GetAll)
		todos.GET("/:id", ctrl.GetByID)
		todos.POST("", ctrl.Create)
		todos.PUT("/:id", ctrl.Update)
		todos.DELETE("/:id", ctrl.Delete)
		todos.PATCH("/:id/complete", ctrl.Complete)
	}

	return &TestHelper{
		DB:         db,
		Router:     router,
		Controller: ctrl,
		Service:    svc,
		Repository: repo,
	}, nil
}

func (h *TestHelper) CleanDatabase() {
	h.DB.Exec("DELETE FROM todos")
}

func (h *TestHelper) CreateTodoRequest(method, url string, body interface{}) (*http.Request, error) {
	var req *http.Request
	var err error

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req = httptest.NewRequest(method, url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	return req, err
}

func (h *TestHelper) ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	h.Router.ServeHTTP(recorder, req)
	return recorder
}

func (h *TestHelper) ParseSuccessResponse(recorder *httptest.ResponseRecorder) (*dto.SuccessResponse, error) {
	var response dto.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	return &response, err
}

func (h *TestHelper) ParseErrorResponse(recorder *httptest.ResponseRecorder) (*dto.ErrorResponse, error) {
	var response dto.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	return &response, err
}
