package repository

import (
	"github.com/pri-laxmi/MailFlow/internal/models"

	"gorm.io/gorm"
)

type CampaignRepository interface {
	Create(campaign *models.Campaign) error

	FindAllByUserID(userID uint) ([]models.Campaign, error)

	FindByIDAndUserID(
		id uint,
		userID uint,
	) (*models.Campaign, error)

	Update(
		id uint,
		userID uint,
		updates map[string]interface{},
	) error

	Delete(
		id uint,
		userID uint,
	) error
}

type campaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) CampaignRepository {
	return &campaignRepository{
		db: db,
	}
}
func (r *campaignRepository) Create(campaign *models.Campaign) error {
	return r.db.Create(campaign).Error
}
func (r *campaignRepository) FindAllByUserID(
	userID uint,
) ([]models.Campaign, error) {

	var campaigns []models.Campaign

	err := r.db.
		Preload("Template").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&campaigns).Error

	return campaigns, err
}
func (r *campaignRepository) FindByIDAndUserID(
	id uint,
	userID uint,
) (*models.Campaign, error) {

	var campaign models.Campaign

	err := r.db.
		Preload("Template").
		Where(
			"id = ? AND user_id = ?",
			id,
			userID,
		).
		First(&campaign).Error

	if err != nil {
		return nil, err
	}

	return &campaign, nil
}
func (r *campaignRepository) Update(
	id uint,
	userID uint,
	updates map[string]interface{},
) error {

	result := r.db.
		Model(&models.Campaign{}).
		Where(
			"id = ? AND user_id = ?",
			id,
			userID,
		).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func (r *campaignRepository) Delete(
	id uint,
	userID uint,
) error {

	result := r.db.
		Where(
			"id = ? AND user_id = ?",
			id,
			userID,
		).
		Delete(&models.Campaign{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
