package workerpool

import "fmt"

type Worker struct {
	id int

	JobPool chan chan Job
	JobChan chan Job

	exit chan struct{}
}

func NewWorker(id int, jobPool chan chan Job) Worker {
	fmt.Printf("create new worker, id: %d\n", id)

	return Worker{
		id:      id,
		JobPool: jobPool,
		JobChan: make(chan Job), // 一個worker 一次只能執行一個job
		exit:    make(chan struct{}),
	}
}

func (ref Worker) Deploy() {

	go func() {
		fmt.Printf("worker %d deployed\n", ref.id)

	LBLJOB:
		for {
			ref.JobPool <- ref.JobChan

			select {
			case job := <-ref.JobChan:
				fmt.Printf("worker %d excuting task %d\n", ref.id, job.Task.Id)

				job.Task.DoSomething()

			case <-ref.exit:
				break LBLJOB
			}
		}

		close(ref.exit)
		close(ref.JobChan)

		fmt.Printf("worker %d dismissed\n", ref.id)
	}()
}

func (ref Worker) Dismiss() {
	go func() {
		ref.exit <- struct{}{}
	}()
}
