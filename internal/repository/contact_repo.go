package repository

import (
	"github.com/pri-laxmi/MailFlow/internal/models"
	"gorm.io/gorm"
)

type ContactRepository interface {
	Create(contact *models.Contact) error
	FindAllByUserID(userID uint) ([]*models.Contact, error)
	FindByIDandUserID(id uint, userID uint) (*models.Contact, error)
	Update(id uint, userID uint, updates map[string]interface{}) (*models.Contact, error)
	Delete(id uint, userID uint) error
}
type contactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{
		db: db,
	}
}

func (r *contactRepository) Create(contact *models.Contact) error {
	return r.db.Create(contact).Error
}
func (r *contactRepository) FindAllByUserID(userID uint) ([]*models.Contact, error) {
	var contact []*models.Contact
	err := r.db.Where("user_id = ?", userID).Find(&contact).Error
	if err != nil {
		return nil, err
	}
	return contact, nil
}
func (r *contactRepository) FindByIDandUserID(id uint, userID uint) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&contact).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}
func (r *contactRepository) Update(
	id uint,
	userID uint,
	updates map[string]interface{},
) (*models.Contact, error) {

	var contact models.Contact

	// Make sure the contact belongs to the authenticated user
	err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		First(&contact).
		Error

	if err != nil {
		return nil, err
	}

	if len(updates) > 0 {
		err = r.db.
			Model(&contact).
			Updates(updates).
			Error

		if err != nil {
			return nil, err
		}
	}

	// Get the updated record
	if err := r.db.First(&contact, contact.ID).Error; err != nil {
		return nil, err
	}

	return &contact, nil

}
func (r *contactRepository) Delete(id uint, userID uint) error {
	result := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Contact{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
