package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/pri-laxmi/MailFlow/internal/dto"
	"github.com/pri-laxmi/MailFlow/internal/service"
)

type ContactHandler struct {
	contactService service.ContactService
}

func NewContactHandler(contactService service.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}
func getUserID(c *gin.Context) (uint, error) {
	value,exists := c.Get("userID")
	if !exists {
		return 0, errors.New("userID not found in context")
	}
	userID, ok := value.(uint)
	if !ok {
		return 0, errors.New("userID is not a valid uint")
	}
	return userID, nil
}
func (h *ContactHandler) Create(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	var req dto.CreateContactRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	contact, err := h.contactService.Create(userID, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create contact",
		})
		return
	}

	c.JSON(http.StatusCreated, contact)
}
func (h *ContactHandler) GetAll(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	contacts, err := h.contactService.GetAll(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch contacts",
		})
		return
	}

	c.JSON(http.StatusOK, contacts)
}
func (h *ContactHandler) GetByID(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	contactID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid contact id",
		})
		return
	}

	contact, err := h.contactService.GetByID(
		userID,
		uint(contactID),
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "contact not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch contact",
		})
		return
	}

	c.JSON(http.StatusOK, contact)
}
func (h *ContactHandler) Update(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	contactID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid contact id",
		})
		return
	}

	var req dto.UpdateContactRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	contact, err := h.contactService.Update(
		userID,
		uint(contactID),
		req,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "contact not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, contact)
}
func (h *ContactHandler) Delete(c *gin.Context) {
	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	contactID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid contact id",
		})
		return
	}

	err = h.contactService.Delete(
		userID,
		uint(contactID),
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "contact not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete contact",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "contact deleted successfully",
	})
}