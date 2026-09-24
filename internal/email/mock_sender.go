package email

import (
	"fmt"
	"log"
	"time"

	"github.com/pri-laxmi/MailFlow/internal/models"
)

type MockSender struct{}

func NewMockSender() Sender {
	return &MockSender{}
}

func (m *MockSender) Send(job *models.Job) error {

	log.Printf(
		"Sending email to %s",
		job.Contact.Email,
	)

	log.Printf(
		"Subject: %s",
		job.Campaign.Template.Subject,
	)

	// Simulate email sending delay
	time.Sleep(1 * time.Second)

	log.Printf(
		"Mock email sent successfully to %s",
		job.Contact.Email,
	)

	fmt.Printf(
		"\n--- MOCK EMAIL ---\nTo: %s\nSubject: %s\nBody: %s\n------------------\n\n",
		job.Contact.Email,
		job.Campaign.Template.Subject,
		job.Campaign.Template.Body,
	)

	return nil
}
