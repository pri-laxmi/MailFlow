package repository

import (
	"github.com/pri-laxmi/MailFlow/internal/models"
	"gorm.io/gorm"
)

type JobLogRepository interface {
	Create(jobLog *models.JobLog) error
	FindByJobID(jobID uint) ([]models.JobLog, error)
}
type jobLogRepository struct {
	db *gorm.DB
}

func NewJobLogRepository(db *gorm.DB) JobLogRepository {
	return &jobLogRepository{
		db: db,
	}
}
func (r *jobLogRepository) Create(jobLog *models.JobLog) error {
	return r.db.Create(jobLog).Error
}
func (r *jobLogRepository) FindByJobID(jobID uint) ([]models.JobLog, error) {
	var jobLogs []models.JobLog
	err := r.db.Where("job_id = ?", jobID).Order("created_at DESC").Find(&jobLogs).Error
	return jobLogs, err
}
