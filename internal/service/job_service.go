package service

import (
	"github.com/pri-laxmi/MailFlow/internal/models"
	"github.com/pri-laxmi/MailFlow/internal/queue"
	"github.com/pri-laxmi/MailFlow/internal/repository"
)

type JobService interface {
	CreateJobsForCampaign(
		campaignID uint,
		contacts []models.Contact,
	) error
}

type jobService struct {
	jobRepo repository.JobRepository
	queue   *queue.Queue
}

func NewJobService(
	jobRepo repository.JobRepository,
	queue *queue.Queue,
) JobService {
	return &jobService{
		jobRepo: jobRepo,
		queue:   queue,
	}
}
func (s *jobService) CreateJobsForCampaign(
	campaignID uint,
	contacts []models.Contact,
) error {

	for _, contact := range contacts {

		job := &models.Job{
			CampaignID: campaignID,
			ContactID:  contact.ID,
			Status:     "pending",
			RetryCount: 0,
		}

		err := s.jobRepo.Create(job)

		if err != nil {
			return err
		}

		s.queue.Push(job)
	}

	return nil
}
