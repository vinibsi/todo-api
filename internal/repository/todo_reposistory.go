package repository

import (
	"github.com/vinibsi/todo-api/internal/entity"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(todo *entity.Todo) error
	GetByID(id uint) (*entity.Todo, error)
	GetAll(limit, offset int) ([]entity.Todo, int64, error)
	Update(todo *entity.Todo) error
	Delete(id uint) error
	GetByCompleted(completed bool, limit, offset int) ([]entity.Todo, int64, error)
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (repo *todoRepository) Create(todo *entity.Todo) error {
	return repo.db.Create(todo).Error
}

func (repo *todoRepository) GetByID(id uint) (*entity.Todo, error) {
	var todo entity.Todo
	err := repo.db.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (repo *todoRepository) GetAll(limit, offset int) ([]entity.Todo, int64, error) {
	var todos []entity.Todo
	var total int64

	// Conta o total de registros
	if err := repo.db.Model(&entity.Todo{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Busca os registros com paginação
	err := repo.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&todos).Error

	return todos, total, err
}

func (repo *todoRepository) Update(todo *entity.Todo) error {
	return repo.db.Save(todo).Error
}

func (repo *todoRepository) Delete(id uint) error {
	return repo.db.Delete(&entity.Todo{}, id).Error
}

func (r *todoRepository) GetByCompleted(completed bool, limit, offset int) ([]entity.Todo, int64, error) {
	var todos []entity.Todo
	var total int64

	query := r.db.Model(&entity.Todo{}).Where("completed = ?", completed)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&todos).Error
	return todos, total, err
}
