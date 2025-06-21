package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vinibsi/todo-api/internal/config"
	"github.com/vinibsi/todo-api/internal/controller"
	"github.com/vinibsi/todo-api/internal/repository"
	"github.com/vinibsi/todo-api/internal/service"
	"github.com/vinibsi/todo-api/pkg/database"
)

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system variables")
	}

	conf := config.Load()

	db, err := database.Connect(conf.DatabaseUrl)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Inicializa camadas
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoController := controller.NewTodoController(todoService)

	// Configura rotas
	router := setupRoutes(todoController)

	// Inicia servidor
	log.Printf("Server running on port %s", conf.Port)
	log.Fatal(router.Run(":" + conf.Port))
}

func setupRoutes(todoController *controller.TodoController) *gin.Engine {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.2", "10.0.0.0/8"})

	api := router.Group("/v1")
	{
		todos := api.Group("/todos")
		{
			todos.GET("", todoController.GetAll)
			todos.GET("/:id", todoController.GetByID)
			todos.POST("", todoController.Create)
			todos.PUT("/:id", todoController.Update)
			todos.DELETE("/:id", todoController.Delete)
			todos.PATCH("/:id/complete", todoController.Complete)
		}
	}

	healthz := router.Group("/healthz")
	{
		healthz.GET("", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"status": "UP",
			})
		})
	}

	return router
}
