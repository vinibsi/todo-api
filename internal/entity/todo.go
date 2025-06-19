package entity

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null;size:255" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Completed   bool           `gorm:"default:false" json:"completed"`
	Priority    string         `gorm:"default:medium;size:20" json:"priority"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
