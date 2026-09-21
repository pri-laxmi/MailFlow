package email

import "github.com/pri-laxmi/MailFlow/internal/models"

type Sender interface {
	Send(job *models.Job) error
}
