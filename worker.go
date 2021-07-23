package goworker

import (
	"time"
)

type ToolQuantity int

const (
	ToolQuantityGood = ToolQuantity(0)
	ToolQuantityBad  = ToolQuantity(1)
)

type IToolMaker interface {
	//Make a tool with origin
	Make(origin string, meta interface{}) interface{}
}

type ITask interface {
	Process(tool interface{}) ToolQuantity
	ToolLabel() string
}

var __tasks = make(chan ITask)
var __num_worker = counter{Value: 0}
var __tickets = make(chan int)

const (
	SSWorkerCommandQuit = 1
)

type worker struct {
	ID            int
	currentOrigin string
}

func (w *worker) goStart() {

	__num_worker.UnsafeLock()
	w.ID = __num_worker.Value
	__num_worker.Value++
	__num_worker.UnsafeUnlock()

	for {

		select {
		case task := <-__tasks:
			toolLabel := task.ToolLabel()
			if toolLabel != "" {
				tool := borrow(toolLabel)
				quantity := task.Process(tool)
				go thankyou(toolLabel, quantity, tool)
			} else {
				task.Process(nil)
			}

		case ticket := <-__tickets:
			if ticket == 1 {
				newWorker := worker{}
				go newWorker.goStart()
			} else {
				__num_worker.Desc()
				return
			}
		default:
			time.Sleep(time.Nanosecond * 10)
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
