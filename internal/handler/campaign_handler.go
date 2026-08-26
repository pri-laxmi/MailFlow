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

type CampaignHandler struct {
	campaignService service.CampaignService
}

func NewCampaignHandler(
	campaignService service.CampaignService,
) *CampaignHandler {
	return &CampaignHandler{
		campaignService: campaignService,
	}
}
func (h *CampaignHandler) Create(c *gin.Context) {

	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	var req dto.CreateCampaignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	campaign, err := h.campaignService.Create(
		userID,
		req,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create campaign",
		})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}
func (h *CampaignHandler) GetAll(c *gin.Context) {

	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	campaigns, err := h.campaignService.GetAll(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch campaigns",
		})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}
func (h *CampaignHandler) GetByID(c *gin.Context) {

	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	campaignID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid campaign id",
		})
		return
	}

	campaign, err := h.campaignService.GetByID(
		userID,
		uint(campaignID),
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "campaign not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch campaign",
		})
		return
	}

	c.JSON(http.StatusOK, campaign)
}
func (h *CampaignHandler) Update(c *gin.Context) {

	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	campaignID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid campaign id",
		})
		return
	}

	var req dto.UpdateCampaignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	campaign, err := h.campaignService.Update(
		userID,
		uint(campaignID),
		req,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "campaign not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, campaign)
}
func (h *CampaignHandler) Delete(c *gin.Context) {

	userID, err := getUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	campaignID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid campaign id",
		})
		return
	}

	err = h.campaignService.Delete(
		userID,
		uint(campaignID),
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "campaign not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete campaign",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "campaign deleted successfully",
	})
}
