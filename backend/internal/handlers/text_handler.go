package handlers

import (
	"net/http"
	"strconv"

	"ancient-text-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type TextHandler struct {
	service *service.AncientTextService
}

func NewTextHandler() *TextHandler {
	return &TextHandler{
		service: service.NewAncientTextService(),
	}
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func (h *TextHandler) Create(c *gin.Context) {
	var req service.CreateTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	text, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    200,
		Message: "success",
		Data:    text,
	})
}

func (h *TextHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	text, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "Text not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data:    text,
	})
}

func (h *TextHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var texts []interface{}
	var total int64
	var err error

	if keyword != "" {
		texts, total, err = h.service.Search(keyword, page, pageSize)
	} else {
		texts, total, err = h.service.List(page, pageSize)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data: PaginatedResponse{
			Items:      texts,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}

func (h *TextHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	text, err := h.service.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data:    text,
	})
}

func (h *TextHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Deleted successfully",
	})
}
