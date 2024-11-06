package workerpool

import (
	"testing"
	"time"
)

// TestWorkerInit 測試worker初始化
func TestWorkerInit(t *testing.T) {
	jobPool := make(chan chan Job, 2)
	worker := NewWorker(1, jobPool)

	if worker.id != 1 {
		t.Errorf("expected worker id to be 1, got %d", worker.id)
	}
	if worker.JobPool != jobPool {
		t.Errorf("expected JobPool to be initialized")
	}
	if worker.JobChan == nil {
		t.Errorf("expected JobChan to be initialized")
	}
	if worker.exit == nil {
		t.Errorf("expected exit channel to be initialized")
	}
}

// // TestJobSchedulerInit 測試job scheduler初始化
func TestJobSchedulerInit(t *testing.T) {
	scheduler := NewJobScheduler(2)

	if scheduler.WorkerCount != 2 {
		t.Errorf("expected WorkerCount to be 2, got %d", scheduler.WorkerCount)
	}
	if scheduler.JobPool == nil {
		t.Errorf("expected JobPool to be initialized")
	}
}

// // TestJobSubmissionAndExecution 測試job提交和執行
func TestJobSubmissionAndExecution(t *testing.T) {

	JobChannel = make(chan Job, 5)

	scheduler := NewJobScheduler(2)
	scheduler.Run()

	// 添加任務並確認執行
	for i := 0; i < 3; i++ {
		task := Task{Id: i}
		job := Job{Task: task}
		JobChannel <- job
	}

	// 等待足夠的時間以確保任務執行
	<-time.After(time.Second * 5)

	scheduler.Shutdown()
}

// TestWorkerExecution 測試worker是否執行Job
func TestWorkerExecution(t *testing.T) {
	jobPool := make(chan chan Job, 1)
	worker := NewWorker(1, jobPool)

	job := Job{
		Task: Task{Id: 1},
	}

	go worker.Deploy()

	select {
	case jobChan := <-jobPool:
		jobChan <- job
	case <-time.After(time.Second):
		t.Errorf("timeout: worker did not accept the job")
	}

	worker.Dismiss()
}

// TestJobSchedulerShutdown 測試JobScheduler關閉
func TestJobSchedulerShutdown(t *testing.T) {
	sched := NewJobScheduler(2)
	sched.Run()

	<-time.After(time.Second) // 給 scheduler 足夠啟動時間
	sched.Shutdown()

	<-time.After(time.Second * CST_WAIT_FOR_DISMISS)
	for _, worker := range sched.Workers {
		select {
		case <-worker.exit:
			// 已關閉

		default:
			t.Errorf("worker %d did not shut down correctly", worker.id)
		}
	}
}

// 測試無任務情況下的 JobScheduler
func TestJobSchedulerNoJob(t *testing.T) {
	JobChannel = make(chan Job, 5)
	scheduler := NewJobScheduler(2)
	scheduler.Run()

	<-time.After(time.Second) // 等待 scheduler 啟動
	scheduler.Shutdown()      // 無任務時關閉 scheduler
}

// 測試邊界情況：JobChannel 滿的情況
func TestJobChannelFull(t *testing.T) {

	JobChannel = make(chan Job, 1)
	sched := NewJobScheduler(2)
	sched.Run()

	// 填滿 JobChannel
	task := Task{Id: 1}
	job := Job{Task: task}
	JobChannel <- job

	select {
	case JobChannel <- job:
		t.Errorf("expected JobChannel to be full, but was able to add job")
	default:
		// 預期行為
	}

	<-time.After(time.Second)

	sched.Shutdown()
}

// TestWorkerDismiss 測試worker關閉
func TestWorkerDismiss(t *testing.T) {
	jobPool := make(chan chan Job, 1)
	worker := NewWorker(1, jobPool)
	go worker.Deploy()

	<-time.After(time.Millisecond * 500)

	worker.Dismiss()

	select {
	case <-worker.exit:
		// 預期行為

	case <-time.After(1 * time.Second):
		t.Errorf("expected worker to dismiss, but it didn't")

	}
}
