package handlers

import (
	"net/http"

	"ancient-text-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AnnotationHandler struct {
	service *service.AnnotationService
}

func NewAnnotationHandler() *AnnotationHandler {
	return &AnnotationHandler{
		service: service.NewAnnotationService(),
	}
}

func (h *AnnotationHandler) Create(c *gin.Context) {
	var req service.CreateAnnotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	ann, err := h.service.Create(&req)
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
		Data:    ann,
	})
}

func (h *AnnotationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	ann, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "Annotation not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data:    ann,
	})
}

func (h *AnnotationHandler) GetByTextID(c *gin.Context) {
	textID := c.Param("textId")
	userID := c.Query("user_id")

	var annotations interface{}
	var err error

	if userID != "" {
		annotations, err = h.service.GetByTextIDAndUser(textID, userID)
	} else {
		annotations, err = h.service.GetByTextID(textID)
	}

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
		Data:    annotations,
	})
}

func (h *AnnotationHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateAnnotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	ann, err := h.service.Update(id, &req)
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
		Data:    ann,
	})
}

func (h *AnnotationHandler) Delete(c *gin.Context) {
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

type BatchSubmitRequest struct {
	TextID      string                        `json:"text_id" binding:"required"`
	UserID      string                        `json:"user_id"`
	Annotations []service.CreateAnnotationRequest `json:"annotations" binding:"required"`
}

func (h *AnnotationHandler) BatchSubmit(c *gin.Context) {
	var req BatchSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	annotations, err := h.service.BatchSubmit(req.TextID, req.UserID, req.Annotations)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    200,
		Message: "Batch submit success",
		Data: gin.H{
			"count":       len(annotations),
			"annotations": annotations,
		},
	})
}
