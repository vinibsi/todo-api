package controller

import (
	"net/http"
	"strconv"

	"github.com/vinibsi/todo-api/internal/dto"
	"github.com/vinibsi/todo-api/internal/service"

	"github.com/gin-gonic/gin"
)

type TodoController struct {
	service service.TodoService
}

func NewTodoController(service service.TodoService) *TodoController {
	return &TodoController{service: service}
}

func (c *TodoController) Create(ctx *gin.Context) {
	var req dto.CreateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid Data",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	todo, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Todo successfully created",
		Data:    todo,
	})
}

func (c *TodoController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	todo, err := c.service.GetByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "todo not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.ErrorResponse{
			Error:   "Failed to get todo",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Data: todo,
	})
}

func (c *TodoController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	todos, err := c.service.GetAll(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Data: todos,
	})
}

func (c *TodoController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req dto.UpdateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid Data",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	todo, err := c.service.Update(uint(id), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "todo not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.ErrorResponse{
			Error:   "Update todo failed",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Todo successfully edited",
		Data:    todo,
	})
}

func (c *TodoController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "todo not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.ErrorResponse{
			Error:   "Delete todo failed",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Todo successfully deleted",
	})
}

func (c *TodoController) Complete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid ID",
			Message: "ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	todo, err := c.service.Complete(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "todo not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, dto.ErrorResponse{
			Error:   "Complete todo failed",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Todo successfully done",
		Data:    todo,
	})
}
