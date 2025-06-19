package dto

import "time"

type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description" binding:"max=1000"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTodoRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string    `json:"description" binding:"omitempty,max=1000"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
	Completed   *bool      `json:"completed"`
}

type TodoResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TodoListResponse struct {
	Data       []TodoResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
