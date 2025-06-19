package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vinibsi/todo-api/internal/config"
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

	log.Println("Database connection established", db)
}
