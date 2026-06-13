package handlers

import (
	"net/http"

	"ancient-text-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type LayoutHandler struct {
	service *service.LayoutConfigService
}

func NewLayoutHandler() *LayoutHandler {
	return &LayoutHandler{
		service: service.NewLayoutConfigService(),
	}
}

func (h *LayoutHandler) GetByTextID(c *gin.Context) {
	textID := c.Param("textId")

	layout, err := h.service.GetByTextID(textID)
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
		Data:    layout,
	})
}

func (h *LayoutHandler) Save(c *gin.Context) {
	textID := c.Param("textId")

	var req service.SaveLayoutConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	layout, err := h.service.Save(textID, &req)
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
		Data:    layout,
	})
}

func (h *LayoutHandler) Delete(c *gin.Context) {
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
