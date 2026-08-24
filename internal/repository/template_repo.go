package repository

import (
	"github.com/pri-laxmi/MailFlow/internal/models"

	"gorm.io/gorm"
)

type TemplateRepository interface {
	Create(template *models.Template) error

	FindAllByUserID(userID uint) ([]models.Template, error)

	FindByIDAndUserID(id uint, userID uint) (*models.Template, error)

	Update(id uint, userID uint, updates map[string]interface{}) error

	Delete(id uint, userID uint) error
}

type templateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepository{
		db: db,
	}
}

func (r *templateRepository) Create(template *models.Template) error {
	return r.db.Create(template).Error
}

func (r *templateRepository) FindAllByUserID(userID uint) ([]models.Template, error) {
	var templates []models.Template

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&templates).Error

	return templates, err
}

func (r *templateRepository) FindByIDAndUserID(
	id uint,
	userID uint,
) (*models.Template, error) {

	var template models.Template

	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		First(&template).Error

	if err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *templateRepository) Update(
	id uint,
	userID uint,
	updates map[string]interface{},
) error {

	result := r.db.
		Model(&models.Template{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *templateRepository) Delete(
	id uint,
	userID uint,
) error {

	result := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Template{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
