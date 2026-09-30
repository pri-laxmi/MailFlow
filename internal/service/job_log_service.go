package service

import (
	"github.com/pri-laxmi/MailFlow/internal/models"
	"github.com/pri-laxmi/MailFlow/internal/repository"
)

type JobLogService interface {
	Log(jobID uint, event string, message string) error
	GetLogs(jobID uint) ([]models.JobLog, error)
}
type jobLogService struct {
	repo repository.JobLogRepository
}

func NewJobLogService(repo repository.JobLogRepository) JobLogService {
	return &jobLogService{
		repo: repo,
	}
}
func (s *jobLogService) Log(jobID uint, event string, message string) error {
	jobLog := &models.JobLog{
		JobID:   jobID,
		Event:   event,
		Message: message,
	}
	return s.repo.Create(jobLog)
}
func (s *jobLogService) GetLogs(jobID uint) ([]models.JobLog, error) {
	return s.repo.FindByJobID(jobID)
}
