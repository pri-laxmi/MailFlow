package service

import (
	"fmt"

	"github.com/pri-laxmi/MailFlow/internal/dto"
	"github.com/pri-laxmi/MailFlow/internal/models"
	"github.com/pri-laxmi/MailFlow/internal/repository"
)

type CampaignService interface {
	Create(
		userID uint,
		req dto.CreateCampaignRequest,
	) (*models.Campaign, error)

	GetAll(
		userID uint,
	) ([]models.Campaign, error)

	GetByID(
		userID uint,
		campaignID uint,
	) (*models.Campaign, error)

	Update(
		userID uint,
		campaignID uint,
		req dto.UpdateCampaignRequest,
	) (*models.Campaign, error)

	Delete(
		userID uint,
		campaignID uint,
	) error
}

type campaignService struct {
	campaignRepo repository.CampaignRepository
}

func NewCampaignService(
	campaignRepo repository.CampaignRepository,
) CampaignService {
	return &campaignService{
		campaignRepo: campaignRepo,
	}
}
func (s *campaignService) Create(
	userID uint,
	req dto.CreateCampaignRequest,
) (*models.Campaign, error) {

	campaign := &models.Campaign{
		UserID:      userID,
		TemplateID:  req.TemplateID,
		Name:        req.Name,
		Status:      "draft",
		ScheduledAt: req.ScheduledAt,
	}

	err := s.campaignRepo.Create(campaign)

	if err != nil {
		return nil, err
	}

	return s.campaignRepo.FindByIDAndUserID(
		campaign.ID,
		userID,
	)
}
func (s *campaignService) GetAll(
	userID uint,
) ([]models.Campaign, error) {

	return s.campaignRepo.FindAllByUserID(userID)
}
func (s *campaignService) GetByID(
	userID uint,
	campaignID uint,
) (*models.Campaign, error) {

	return s.campaignRepo.FindByIDAndUserID(
		campaignID,
		userID,
	)
}
func (s *campaignService) Update(
	userID uint,
	campaignID uint,
	req dto.UpdateCampaignRequest,
) (*models.Campaign, error) {

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.TemplateID != nil {
		updates["template_id"] = *req.TemplateID
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if req.ScheduledAt != nil {
		updates["scheduled_at"] = *req.ScheduledAt
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	err := s.campaignRepo.Update(
		campaignID,
		userID,
		updates,
	)

	if err != nil {
		return nil, err
	}

	return s.campaignRepo.FindByIDAndUserID(
		campaignID,
		userID,
	)
}
func (s *campaignService) Delete(
	userID uint,
	campaignID uint,
) error {

	return s.campaignRepo.Delete(
		campaignID,
		userID,
	)
}
