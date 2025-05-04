package Tcontrol

import "sync"

type TaskContro struct {
	wg   sync.WaitGroup
	done chan struct{}
}

func NewTaskContro() *TaskContro {
	return &TaskContro{
		done: make(chan struct{}),
	}
}

func (task *TaskContro) Begin() {
	task.wg.Add(1)
}

func (task *TaskContro) Done() {
	task.wg.Done()
}

func (task *TaskContro) Wait() {
	task.wg.Wait()
	close(task.done) // optional signal to others
}

func (task *TaskContro) DoneChan() <-chan struct{} {
	return task.done
}
