package service

//A user can only create, view, update, or delete their own contacts.
import (
	"errors"

	"github.com/pri-laxmi/MailFlow/internal/dto"
	"github.com/pri-laxmi/MailFlow/internal/models"
	"github.com/pri-laxmi/MailFlow/internal/repository"
)

type ContactService interface {
	Create(userID uint, req dto.CreateContactRequest) (*dto.ContactResponse, error)
	GetAll(userID uint) ([]dto.ContactResponse, error)
	GetByID(userID uint, contactID uint) (*dto.ContactResponse, error)
	Update(userID uint, contactID uint, req dto.UpdateContactRequest) (*dto.ContactResponse, error)
	Delete(userID uint, contactID uint) error
}
type contactService struct {
	ContactRepo repository.ContactRepository
}

func NewContactService(contactRepo repository.ContactRepository) ContactService {
	return &contactService{
		ContactRepo: contactRepo,
	}
}
func (s *contactService) Create(userID uint, req dto.CreateContactRequest) (*dto.ContactResponse, error) {
	contact := &models.Contact{
		Name:   req.Name,
		Email:  req.Email,
		UserID: userID,
	}
	err := s.ContactRepo.Create(contact)
	if err != nil {
		return nil, err
	}
	return &dto.ContactResponse{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}, nil
}
func (s *contactService) GetAll(userID uint) ([]dto.ContactResponse, error) {
	contacts, err := s.ContactRepo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	responses := make([]dto.ContactResponse, 0, len(contacts))
	for _, contact := range contacts {
		responses = append(responses, dto.ContactResponse{
			ID:        contact.ID,
			Name:      contact.Name,
			Email:     contact.Email,
			CreatedAt: contact.CreatedAt,
			UpdatedAt: contact.UpdatedAt,
		})
	}
	return responses, nil
}
func (s *contactService) GetByID(userID uint, contactID uint) (*dto.ContactResponse, error) {
	contact, err := s.ContactRepo.FindByIDandUserID(contactID, userID)
	if err != nil {
		return nil, err
	}
	return &dto.ContactResponse{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}, nil
}
func (s *contactService) Update(
	userID uint,
	contactID uint,
	req dto.UpdateContactRequest,
) (*dto.ContactResponse, error) {

	if req.Name == nil && req.Email == nil {
		return nil, errors.New("at least one field is required")
	}
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	contact, err := s.ContactRepo.Update(
		contactID,
		userID,
		updates,
	)

	if err != nil {
		return nil, err
	}

	return &dto.ContactResponse{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}, nil
}
func (s *contactService) Delete(
	userID uint,
	contactID uint,
) error {

	return s.ContactRepo.Delete(
		contactID,
		userID,
	)
}
