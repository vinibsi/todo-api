package database

import (
	"github.com/vinibsi/todo-api/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(databaseUrl string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Detecta o tipo de banco de dados a partir da URL
	if databaseUrl == ":memory:" || databaseUrl == "file::memory:?cache=shared" {
		// SQLite em memória para testes
		db, err = gorm.Open(sqlite.Open(databaseUrl), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		// PostgreSQL para produção
		db, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		return nil, err
	}

	// Auto-migração
	if err := db.AutoMigrate(&entity.Todo{}); err != nil {
		return nil, err
	}

	return db, nil

}

func ConnectTest() (*gorm.DB, error) {
	return Connect(":memory:")
}
