package worker

import (
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
	queue   *queue.Queue
	jobRepo repository.JobRepository
	sender  email.Sender
	workers int
	wg      sync.WaitGroup //wait for all the workers to finish during shutdown
	jobLogService service.JobLogService
}

func NewWorkerPool(queue *queue.Queue, jobRepo repository.JobRepository, sender email.Sender, workers int) *WorkerPool {
	return &WorkerPool{
		queue:   queue,
		jobRepo: jobRepo,
		sender:  sender,
		workers: workers,
		jobLogService: service.NewJobLogService(),
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
	err := p.jobRepo.UpdateStatus(
		job.ID,
		"processing",
		"",
	)

	if err != nil {
		return err
	}
	//load complete job info
	fullJob, err := p.jobRepo.FindByID(job.ID)
	if err != nil {
		return err
	}
	err = p.sender.Send(fullJob)
	if err == nil {
		p.jobRepo.UpdateStatus(
			job.ID,
			"completed",
			"",
		)
		return err
	}

	//mark as failed
	if err != nil {
		p.jobRepo.UpdateStatus(
			job.ID,
			"failed",
			err.Error(),
		)
		return err
	}

	/*if err != nil {
		return err
	}*/
	fullJob.RetryCount++
	if fullJob.RetryCount > MaxRetries {
		p.jobRepo.UpdateStatus(
			job.ID,
			"failed",
			"max retries exceeded",
		)
		return err
	}
	err = p.jobRepo.UpdateRetryCount(
		job.ID,
		fullJob.RetryCount,
		err.Error(),
	)
	if err != nil {
		return err
	}
	delay := retryDelay(fullJob.RetryCount)
	log.Printf(
		"Retrying job %d after %v seconds (retry count: %d)")
	time.Sleep(delay)

	return nil
}
func retryDelay(retryCount int) time.Duration {
	return time.Duration(1<<uint(retryCount-1)) * time.Second
}
