package goworker_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tapvanvn/goworker"
)

type testTask struct {
	Value int
}

func (t *testTask) Process() {
	fmt.Println("Process task", t.Value)
}

func TestWorker(t *testing.T) {
	goworker.OrganizeWorker(5)
	for i := 0; i < 10; i++ {
		task := testTask{
			Value: i,
		}
		goworker.AddTask(&task)
	}
	time.Sleep(time.Second)
	goworker.OrganizeWorker(0)
}
