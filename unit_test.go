package goworker_test

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
	"time"

	crypto_rand "crypto/rand"

	"github.com/tapvanvn/goworker/v2"
)

var randArrayCase string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

//GenVerifyCode generate a verify code
func genName(length int) string {

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	var code string = ""
	var arrayLen = len(randArrayCase)
	for i := 0; i < length; i++ {
		code += string(randArrayCase[rand.Intn(arrayLen)])
	}

	return code
}

//MARK: testTask
type testTask struct {
	Value int
}

func (t *testTask) Process(tool interface{}) goworker.ToolQuantity {
	fmt.Println("Process task", t.Value)
	return goworker.ToolQuantityGood

}
func (t *testTask) ToolLabel() string {
	return ""
}

//MARK: testTaskWithTool
type testTool struct {
	Name string
}
type testToolMaker struct {
	maked int
}

func (t *testToolMaker) Make(origin string, meta interface{}) interface{} {
	t.maked++
	return &testTool{
		Name: fmt.Sprint(t.maked),
	}
}

type TaskWithTool struct {
	Value int
}

func (t *TaskWithTool) Process(tool interface{}) goworker.ToolQuantity {
	if testtool, ok := tool.(*testTool); ok {
		fmt.Println("Process task", t.Value, "with tool", testtool.Name, "success")
	} else {
		fmt.Println("Process task with error tool", t.Value)
	}
	return goworker.ToolQuantityGood
}
func (t *TaskWithTool) ToolLabel() string {
	return "test"
}

//BAD
type TaskWithBadTool struct {
	Value int
}

func (t *TaskWithBadTool) Process(tool interface{}) goworker.ToolQuantity {

	testtool := tool.(*testTool)

	if rand.Intn(3) == 1 {

		fmt.Println("Process task", t.Value, "with tool", testtool.Name, "success")
		return goworker.ToolQuantityGood
	}
	fmt.Println("Process task", t.Value, "with tool", testtool.Name, "false")
	return goworker.ToolQuantityBad
}

func (t *TaskWithBadTool) ToolLabel() string {
	return "test"
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

func TestWorkerWithTool(t *testing.T) {
	toolMaker := testToolMaker{}
	goworker.AddTool("test", &toolMaker)

	goworker.OrganizeWorker(5)

	for i := 0; i < 5; i++ {
		task := TaskWithTool{
			Value: i,
		}
		goworker.AddTask(&task)
	}
	time.Sleep(time.Second * 2)

	for i := 0; i < 5; i++ {
		task := TaskWithTool{
			Value: i,
		}
		goworker.AddTask(&task)
	}

	goworker.OrganizeWorker(0)
}

func TestWorkerWithToolControll(t *testing.T) {
	toolMaker := testToolMaker{}
	goworker.AddToolWithControl("test", &toolMaker, 10)

	goworker.OrganizeWorker(5)

	for i := 0; i < 5; i++ {
		task := TaskWithTool{
			Value: i,
		}
		goworker.AddTask(&task)
	}
	time.Sleep(time.Second * 2)

	for i := 0; i < 5; i++ {
		task := TaskWithTool{
			Value: i,
		}
		goworker.AddTask(&task)
	}

	goworker.OrganizeWorker(0)
}

func TestWorkerWithBadToolControll(t *testing.T) {

	toolMaker := testToolMaker{}
	goworker.AddToolWithControl("test", &toolMaker, 1)
	goworker.OrganizeWorker(1)

	for i := 0; i < 5; i++ {
		var task goworker.ITask = nil

		if rand.Intn(2) == 1 {
			fmt.Println("bad", i)
			task = &TaskWithBadTool{
				Value: i,
			}
		} else {
			fmt.Println("good", i)
			task = &TaskWithTool{
				Value: i,
			}
		}

		goworker.AddTask(task)
	}
	time.Sleep(time.Second)
	time.Sleep(time.Second * 3)
	//panic("")
}
