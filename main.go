package main

import (
	"fmt"
	"time"

	workerPool "sw.worker.pool/workerpool"
)

func main() {

	workerPool.JobChannel = make(chan workerPool.Job, 100)

	sched := workerPool.NewJobScheduler(10)
	sched.Run()

	go addJob1()
	<-time.After(time.Second * 10)
	go addJob2()

	<-time.After(time.Second * 20)
	fmt.Println("shutdown")
	sched.Shutdown()

	<-time.After(time.Hour)
}

func addJob1() {
	<-time.After(time.Second * 1)
	fmt.Println("add jobs 1")
	for i := 0; i < 2; i++ {
		task := workerPool.Task{Id: i}
		job := workerPool.Job{Task: task}
		workerPool.JobChannel <- job
	}
}

func addJob2() {
	fmt.Println("add jobs 2")
	for i := 2; i < 7; i++ {
		task := workerPool.Task{Id: i}
		job := workerPool.Job{Task: task}
		workerPool.JobChannel <- job
		<-time.After(time.Millisecond * 250)
	}
}

// const (
// 	CST_WAIT_FOR_DISMISS = 2
// )

// var JobChannel chan Job

// type Job struct {
// 	Task Task
// }

// type Task struct {
// 	Id int
// }

// func (ref Task) DoSomething() {

// 	// 模擬task 執行時間
// 	<-time.After(time.Second * 1)

// 	fmt.Printf("%d task done\n", ref.Id)
// }

// type JobScheduler struct {
// 	WorkerCount int
// 	Workers     []*Worker

// 	JobPool chan chan Job
// }

// func NewJobScheduler(workerCount int) *JobScheduler {

// 	return &JobScheduler{
// 		WorkerCount: workerCount,
// 		JobPool:     make(chan chan Job, workerCount),
// 	}
// }

// func (ref *JobScheduler) run() {

// 	for job := range JobChannel {

// 		go func(job Job) {
// 			jobChan := <-ref.JobPool
// 			jobChan <- job

// 		}(job)
// 	}
// }

// func (ref *JobScheduler) Run() {

// 	workers := make([]*Worker, ref.WorkerCount)

// 	for i := 0; i < ref.WorkerCount; i++ {
// 		worker := NewWorker(i, ref.JobPool)
// 		worker.Deploy()

// 		workers[i] = &worker
// 	}

// 	ref.Workers = workers

// 	go ref.run()
// }

// func (ref *JobScheduler) Shutdown() {
// 	for _, worker := range ref.Workers {
// 		worker.Dismiss()
// 	}

// 	<-time.After(time.Second * CST_WAIT_FOR_DISMISS)

// 	close(ref.JobPool)
// }

// type Worker struct {
// 	id int

// 	JobPool chan chan Job
// 	JobChan chan Job

// 	exit chan struct{}
// }

// func NewWorker(id int, jobPool chan chan Job) Worker {
// 	fmt.Printf("create new worker, id: %d\n", id)

// 	return Worker{
// 		id:      id,
// 		JobPool: jobPool,
// 		JobChan: make(chan Job), // 一個worker 一次只能執行一個job
// 		exit:    make(chan struct{}),
// 	}
// }

// func (ref Worker) Deploy() {

// 	go func() {
// 		fmt.Printf("worker %d deployed\n", ref.id)

// 	LBLJOB:
// 		for {
// 			ref.JobPool <- ref.JobChan

// 			select {
// 			case job := <-ref.JobChan:
// 				job.Task.DoSomething()

// 			case <-ref.exit:
// 				break LBLJOB
// 			}
// 		}

// 		close(ref.exit)
// 		close(ref.JobChan)

// 		fmt.Printf("worker %d dismissed\n", ref.id)
// 	}()
// }

// func (ref Worker) Dismiss() {
// 	go func() {
// 		ref.exit <- struct{}{}
// 	}()
// }
