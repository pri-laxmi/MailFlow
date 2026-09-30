package worker

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pri-laxmi/MailFlow/internal/email"
	"github.com/pri-laxmi/MailFlow/internal/models"
	"github.com/pri-laxmi/MailFlow/internal/queue"
	"github.com/pri-laxmi/MailFlow/internal/repository"
	"github.com/pri-laxmi/MailFlow/internal/service"
)

type WorkerPool struct {
	queue         *queue.Queue
	jobRepo       repository.JobRepository
	sender        email.Sender
	workers       int
	wg            sync.WaitGroup //wait for all the workers to finish during shutdown
	jobLogService service.JobLogService
}

func NewWorkerPool(queue *queue.Queue, jobRepo repository.JobRepository, jobLogRepo repository.JobLogRepository, sender email.Sender, workers int) *WorkerPool {
	return &WorkerPool{
		queue:         queue,
		jobRepo:       jobRepo,
		sender:        sender,
		workers:       workers,
		jobLogService: service.NewJobLogService(jobLogRepo),
	}
}
func (p *WorkerPool) Start() {
	for i := 1; i <= p.workers; i++ {
		p.wg.Add(1)

		go p.worker(i)
	}

	log.Printf("worker pool started with %d workers", p.workers)
}
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	log.Printf("worker %d started", id)

	for {
		job := p.queue.Pop()

		log.Printf(
			"worker %d picked job %d",
			id,
			job.ID,
		)

		if err := p.processJob(job); err != nil {
			log.Printf(
				"worker %d failed job %d: %v",
				id,
				job.ID,
				err,
			)

			continue
		}

		log.Printf(
			"worker %d completed job %d",
			id,
			job.ID,
		)
	}
}

// change this
func (p *WorkerPool) processJob(job *models.Job) error {

	// --------------------------------------------------
	// 1. Load complete job information
	// --------------------------------------------------

	fullJob, err := p.jobRepo.FindByID(job.ID)

	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 2. Mark job as processing
	// --------------------------------------------------

	err = p.jobRepo.UpdateStatus(
		job.ID,
		"processing",
		"",
	)

	if err != nil {
		return err
	}

	// Log processing started
	err = p.jobLogService.Log(
		job.ID,
		service.ProcessingStarted,
		"Worker started processing job",
	)

	if err != nil {
		log.Printf(
			"failed to create processing log for job %d: %v",
			job.ID,
			err,
		)
	}

	// --------------------------------------------------
	// 3. Try sending the email
	// --------------------------------------------------

	err = p.sender.Send(fullJob)

	// --------------------------------------------------
	// 4. Email sent successfully
	// --------------------------------------------------

	if err == nil {

		// Log email sent
		logErr := p.jobLogService.Log(
			job.ID,
			service.EmailSent,
			"Email sent successfully",
		)

		if logErr != nil {
			log.Printf(
				"failed to create email sent log for job %d: %v",
				job.ID,
				logErr,
			)
		}

		// Mark job as completed
		err = p.jobRepo.UpdateStatus(
			job.ID,
			"completed",
			"",
		)

		if err != nil {
			return err
		}

		// Log job completed
		logErr = p.jobLogService.Log(
			job.ID,
			service.JobCompleted,
			"Job completed successfully",
		)

		if logErr != nil {
			log.Printf(
				"failed to create completion log for job %d: %v",
				job.ID,
				logErr,
			)
		}

		return nil
	}

	// --------------------------------------------------
	// 5. Email failed
	// --------------------------------------------------

	log.Printf(
		"job %d email sending failed: %v",
		job.ID,
		err,
	)

	// Log email failure
	logErr := p.jobLogService.Log(
		job.ID,
		service.EmailFailed,
		err.Error(),
	)

	if logErr != nil {
		log.Printf(
			"failed to create failure log for job %d: %v",
			job.ID,
			logErr,
		)
	}

	// --------------------------------------------------
	// 6. Increase retry count
	// --------------------------------------------------

	fullJob.RetryCount++

	// --------------------------------------------------
	// 7. Maximum retries reached
	// --------------------------------------------------

	if fullJob.RetryCount >= MaxRetries {

		log.Printf(
			"job %d permanently failed after %d attempts",
			job.ID,
			fullJob.RetryCount,
		)

		// Mark job as failed
		updateErr := p.jobRepo.UpdateStatus(
			job.ID,
			"failed",
			err.Error(),
		)

		if updateErr != nil {
			return updateErr
		}

		// Log permanent failure
		logErr = p.jobLogService.Log(
			job.ID,
			service.JobFailed,
			"Maximum retry attempts reached",
		)

		if logErr != nil {
			log.Printf(
				"failed to create permanent failure log for job %d: %v",
				job.ID,
				logErr,
			)
		}

		return err
	}

	// --------------------------------------------------
	// 8. Save retry count
	// --------------------------------------------------

	err = p.jobRepo.UpdateRetryCount(
		job.ID,
		fullJob.RetryCount,
		err.Error(),
	)

	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 9. Calculate exponential backoff
	// --------------------------------------------------

	delay := retryDelay(fullJob.RetryCount)

	log.Printf(
		"job %d retrying in %v",
		job.ID,
		delay,
	)

	// Log retry scheduled
	logErr = p.jobLogService.Log(
		job.ID,
		service.RetryScheduled,
		fmt.Sprintf(
			"Retry #%d scheduled after %v",
			fullJob.RetryCount,
			delay,
		),
	)

	if logErr != nil {
		log.Printf(
			"failed to create retry log for job %d: %v",
			job.ID,
			logErr,
		)
	}

	// --------------------------------------------------
	// 10. Wait before retrying
	// --------------------------------------------------

	time.Sleep(delay)

	// --------------------------------------------------
	// 11. Retry the job
	// --------------------------------------------------

	return p.processJob(job)
}
func retryDelay(retryCount int) time.Duration {
	return time.Duration(1<<uint(retryCount-1)) * time.Second
}
