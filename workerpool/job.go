package workerpool

var JobChannel chan Job

type Job struct {
	Task Task
}
