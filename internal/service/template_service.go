package service

import (
	"fmt"

	"github.com/pri-laxmi/MailFlow/internal/dto"
	"github.com/pri-laxmi/MailFlow/internal/models"
	"github.com/pri-laxmi/MailFlow/internal/repository"
)

type TemplateService interface {
	Create(userID uint, req dto.CreateTemplateRequest) (*models.Template, error)

	GetAll(userID uint) ([]models.Template, error)

	GetByID(userID uint, templateID uint) (*models.Template, error)

	Update(
		userID uint,
		templateID uint,
		req dto.UpdateTemplateRequest,
	) (*models.Template, error)

	Delete(userID uint, templateID uint) error
}

type templateService struct {
	templateRepo repository.TemplateRepository
}

func NewTemplateService(
	templateRepo repository.TemplateRepository,
) TemplateService {
	return &templateService{
		templateRepo: templateRepo,
	}
}
func (s *templateService) Create(
	userID uint,
	req dto.CreateTemplateRequest,
) (*models.Template, error) {

	template := &models.Template{
		UserId:  userID,
		Name:    req.Name,
		Subject: req.Subject,
		Body:    req.Body,
	}

	err := s.templateRepo.Create(template)

	if err != nil {
		return nil, err
	}

	return template, nil
}
func (s *templateService) GetAll(
	userID uint,
) ([]models.Template, error) {

	return s.templateRepo.FindAllByUserID(userID)
}
func (s *templateService) GetByID(
	userID uint,
	templateID uint,
) (*models.Template, error) {

	return s.templateRepo.FindByIDAndUserID(
		templateID,
		userID,
	)
}
func (s *templateService) Update(
	userID uint,
	templateID uint,
	req dto.UpdateTemplateRequest,
) (*models.Template, error) {

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}

	if req.Body != nil {
		updates["body"] = *req.Body
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	err := s.templateRepo.Update(
		templateID,
		userID,
		updates,
	)

	if err != nil {
		return nil, err
	}

	return s.templateRepo.FindByIDAndUserID(
		templateID,
		userID,
	)
}
func (s *templateService) Delete(
	userID uint,
	templateID uint,
) error {

	return s.templateRepo.Delete(
		templateID,
		userID,
	)
}
