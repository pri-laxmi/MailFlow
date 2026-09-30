package repository

import (
	"github.com/pri-laxmi/MailFlow/internal/models"

	"gorm.io/gorm"
)

type JobRepository interface {
	Create(job *models.Job) error

	FindByID(id uint) (*models.Job, error)

	FindPending() ([]models.Job, error)

	FindAllByCampaignID(campaignID uint) ([]models.Job, error)

	UpdateStatus(
		id uint,
		status string,
		errorMessage string,
	) error
	UpdateRetryCount(id uint, retryCount int, errorMessage string) error
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{
		db: db,
	}
}
func (r *jobRepository) Create(job *models.Job) error {
	return r.db.Create(job).Error
}
func (r *jobRepository) FindByID(
	id uint,
) (*models.Job, error) {

	var job models.Job

	err := r.db.
		Preload("Campaign").
		Preload("Campaign.Template").
		Preload("Contact").
		First(&job, id).Error

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *jobRepository) FindPending() ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.
		Where("status = ?", "pending").
		Order("created_at ASC").
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) FindAllByCampaignID(
	campaignID uint,
) ([]models.Job, error) {

	var jobs []models.Job

	err := r.db.
		Where("campaign_id = ?", campaignID).
		Order("created_at ASC").
		Find(&jobs).Error

	return jobs, err
}
func (r *jobRepository) UpdateStatus(
	id uint,
	status string,
	errorMessage string,
) error {

	return r.db.
		Model(&models.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        status,
			"error_message": errorMessage,
		}).Error
}
func (r *jobRepository) UpdateRetryCount(
	id uint,
	retryCount int,
	errorMessage string,
) error {
	return r.db.
		Model(&models.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"retry_count":   retryCount,
			"error_message": errorMessage,
		}).Error
}
