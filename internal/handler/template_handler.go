package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pri-laxmi/MailFlow/internal/dto"
	"github.com/pri-laxmi/MailFlow/internal/service"
	"gorm.io/gorm"
)

type TemplateHandler struct {
	templateService service.TemplateService
}

func NewTemplateHandler(
	templateService service.TemplateService,
) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
	}
}
func (h *TemplateHandler) Create(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	var req dto.CreateTemplateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	template, err := h.templateService.Create(userID, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create template",
		})
		return
	}

	c.JSON(http.StatusCreated, template)
}
func (h *TemplateHandler) GetByID(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	templateID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid template id",
		})
		return
	}

	template, err := h.templateService.GetByID(
		userID,
		uint(templateID),
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "template not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch template",
		})
		return
	}

	c.JSON(http.StatusOK, template)
}
func (h *TemplateHandler) GetAll(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	templates, err := h.templateService.GetAll(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch templates",
		})
		return
	}

	c.JSON(http.StatusOK, templates)
}
func (h *TemplateHandler) Update(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	templateID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid template id",
		})
		return
	}

	var req dto.UpdateTemplateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	template, err := h.templateService.Update(
		userID,
		uint(templateID),
		req,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "template not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, template)
}
func (h *TemplateHandler) Delete(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	templateID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid template id",
		})
		return
	}

	err = h.templateService.Delete(
		userID,
		uint(templateID),
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "template not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete template",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "template deleted successfully",
	})
}
