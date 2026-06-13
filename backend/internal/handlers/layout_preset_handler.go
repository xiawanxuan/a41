package handlers

import (
	"net/http"

	"ancient-text-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type LayoutPresetHandler struct {
	service *service.LayoutPresetService
}

func NewLayoutPresetHandler() *LayoutPresetHandler {
	return &LayoutPresetHandler{
		service: service.NewLayoutPresetService(),
	}
}

func getUserID(c *gin.Context) string {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "default-user"
	}
	return userID
}

func (h *LayoutPresetHandler) Create(c *gin.Context) {
	userID := getUserID(c)

	var req service.CreateLayoutPresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	preset, err := h.service.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    preset,
	})
}

func (h *LayoutPresetHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	preset, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Preset not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    preset,
	})
}

func (h *LayoutPresetHandler) ListByUser(c *gin.Context) {
	userID := getUserID(c)

	presets, err := h.service.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    presets,
	})
}

func (h *LayoutPresetHandler) GetDefault(c *gin.Context) {
	userID := getUserID(c)

	preset, err := h.service.GetDefault(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    preset,
	})
}

func (h *LayoutPresetHandler) Update(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")

	var req service.UpdateLayoutPresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	preset, err := h.service.Update(id, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    preset,
	})
}

func (h *LayoutPresetHandler) SetDefault(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")

	err := h.service.SetDefault(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

func (h *LayoutPresetHandler) Delete(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")

	err := h.service.Delete(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}
