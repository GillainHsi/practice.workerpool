package workerpool

import (
	"time"
)

type JobScheduler struct {
	WorkerCount int
	Workers     []*Worker

	JobPool chan chan Job
}

func NewJobScheduler(workerCount int) *JobScheduler {

	return &JobScheduler{
		WorkerCount: workerCount,
		JobPool:     make(chan chan Job, workerCount),
	}
}

func (ref *JobScheduler) run() {

	for job := range JobChannel {

		go func(job Job) {
			jobChan := <-ref.JobPool
			jobChan <- job

		}(job)
	}
}

func (ref *JobScheduler) Run() {

	workers := make([]*Worker, ref.WorkerCount)

	for i := 0; i < ref.WorkerCount; i++ {
		worker := NewWorker(i, ref.JobPool)
		worker.Deploy()

		workers[i] = &worker
	}

	ref.Workers = workers

	go ref.run()
}

func (ref *JobScheduler) Shutdown() {
	for _, worker := range ref.Workers {
		worker.Dismiss()
	}

	<-time.After(time.Second * CST_WAIT_FOR_DISMISS)

	close(ref.JobPool)
}
