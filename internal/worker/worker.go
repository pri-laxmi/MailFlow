package worker

import (
	"log"
	"sync"

	"github.com/pri-laxmi/MailFlow/internal/email"
	"github.com/pri-laxmi/MailFlow/internal/models"
	"github.com/pri-laxmi/MailFlow/internal/queue"
	"github.com/pri-laxmi/MailFlow/internal/repository"
)

type WorkerPool struct {
	queue   *queue.Queue
	jobRepo repository.JobRepository
	sender  email.Sender
	workers int
	wg      sync.WaitGroup //wait for all the workers to finish during shutdown
}

func NewWorkerPool(queue *queue.Queue, jobRepo repository.JobRepository, sender email.Sender, workers int) *WorkerPool {
	return &WorkerPool{
		queue:   queue,
		jobRepo: jobRepo,
		sender:  sender,
		workers: workers,
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
	if err != nil {
		p.jobRepo.UpdateStatus(
			job.ID,
			"failed",
			err.Error(),
		)
		return err
	}

	//mark as completed
	err = p.jobRepo.UpdateStatus(
		job.ID,
		"completed",
		"",
	)

	if err != nil {
		return err
	}

	return nil
}
