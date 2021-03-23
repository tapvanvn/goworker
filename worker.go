package goworker

import (
	"sync"
)

type ITask interface {
	Process()
}

var __tasks = make(chan ITask)
var __num_worker = counter{Value: 0}
var __tickets = make(chan int)

type counter struct {
	mu    sync.Mutex
	Value int
}

func (c *counter) Inc() {
	c.mu.Lock()
	c.Value++
	c.mu.Unlock()
}

func (c *counter) Desc() {
	c.mu.Lock()
	c.Value--
	c.mu.Unlock()
}
func (c *counter) UnsafeLock() {
	c.mu.Lock()
}
func (c *counter) UnsafeUnlock() {
	c.mu.Unlock()
}

const (
	SSWorkerCommandQuit = 1
)

type worker struct {
	ID int
}

func (w *worker) goStart() {

	__num_worker.UnsafeLock()
	w.ID = __num_worker.Value
	__num_worker.Value++
	__num_worker.UnsafeUnlock()

	for {

		select {
		case task := <-__tasks:

			task.Process()

		case ticket := <-__tickets:
			if ticket == 1 {
				newWorker := worker{}
				go newWorker.goStart()
			} else {
				__num_worker.Desc()
				return
			}
		}
	}
}

func OrganizeWorker(numWorker int) {

	if __num_worker.Value == 0 {
		for i := 0; i < numWorker; i++ {
			newWorker := worker{}
			go newWorker.goStart()
		}
		return
	}

	__num_worker.UnsafeLock()
	curr := __num_worker.Value

	if numWorker > curr {
		for i := curr; i < numWorker; i++ {
			__tickets <- 1
		}
	} else if numWorker < curr {
		for i := curr; i > numWorker; i-- {
			__tickets <- -1
		}
	}
	__num_worker.UnsafeUnlock()
}

func AddTask(task ITask) {
	__tasks <- task
}
