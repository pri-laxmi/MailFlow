package queue

import "github.com/pri-laxmi/MailFlow/internal/models"

type Queue struct {
	Jobs chan *models.Job
}

func NewQueue(bufferSize int) *Queue {
	return &Queue{
		Jobs: make(chan *models.Job, bufferSize),
	}
}

func (q *Queue) Push(job *models.Job) {
	q.Jobs <- job
}

func (q *Queue) Pop() *models.Job {
	return <-q.Jobs
}
