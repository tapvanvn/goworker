package goworker

import (
	"math"

	"github.com/tapvanvn/gowrandom"
)

type QuantityReporter func(origin string, meta interface{}, failCount int)
type ToolHandle struct {
	origin string
	tool   interface{}
}
type ToolManager struct {
	BlackSmith      IToolMaker
	numMade         int
	stack           chan *ToolHandle
	max             int
	wrand           *gowrandom.WRandom
	wrandIndex      map[int]string
	quantityControl map[string]int
	originMeta      map[string]interface{}
	reporter        QuantityReporter
}

func (m *ToolManager) pickOrigin() (string, interface{}) {
	index := m.wrand.Pick()
	origin := m.wrandIndex[index]
	return origin, m.originMeta[origin]
}

var __tools = map[string]*ToolManager{}

func AddTool(label string, toolMaker IToolMaker) {

	AddToolWithControl(label, toolMaker, math.MaxInt32)
}

func AddToolWithControl(label string, toolMaker IToolMaker, maxTool int) {

	if manager, ok := __tools[label]; !ok {

		m := &ToolManager{
			BlackSmith:      toolMaker,
			stack:           make(chan *ToolHandle, 1),
			numMade:         1,
			max:             maxTool,
			wrand:           gowrandom.MakeWRandom(1),
			wrandIndex:      map[int]string{},
			quantityControl: map[string]int{},
			originMeta:      map[string]interface{}{},
			reporter:        nil,
		}
		m.wrand.SetWeight(0, 100)
		m.wrandIndex[0] = "default"
		origin, meta := m.pickOrigin()
		toolHandle := &ToolHandle{
			origin: origin,
			tool:   m.BlackSmith.Make(origin, meta),
		}
		m.stack <- toolHandle
		__tools[label] = m

	} else {

		manager.BlackSmith = toolMaker
	}
}

//borrow a tool to do job
func borrow(label string) *ToolHandle {

	if manager, ok := __tools[label]; ok {
		if len(manager.stack) == 0 && manager.numMade < manager.max {
			manager.numMade++
			origin, meta := manager.pickOrigin()
			toolHandle := &ToolHandle{
				origin: origin,
				tool:   manager.BlackSmith.Make(origin, meta),
			}
			return toolHandle
		}

		return <-manager.stack
	}
	return nil
}

//release tool after using
func thankyou(label string, quantity ToolQuantity, tool *ToolHandle) {
	if quantity == ToolQuantityBad {

		manager, _ := __tools[label]
		manager.quantityControl[tool.origin]++

		origin, meta := manager.pickOrigin()
		toolHandle := &ToolHandle{
			origin: origin,
			tool:   manager.BlackSmith.Make(origin, meta),
		}
		manager.stack <- toolHandle
		if manager.reporter != nil {
			reportMeta := manager.originMeta[tool.origin]
			go manager.reporter(tool.origin, reportMeta, manager.quantityControl[tool.origin])
		}
	} else {
		if manager, ok := __tools[label]; ok {
			manager.stack <- tool
		}
	}
}

func AddOrigin(label string, origin string, meta interface{}, randomWeight int) {
	if manager, ok := __tools[label]; ok {
		for index, ori := range manager.wrandIndex {
			if ori == origin {
				manager.wrand.SetWeight(index, uint(randomWeight))
				manager.originMeta[ori] = meta
				return
			}
		}
		newIndex := manager.wrand.AddElement(uint(randomWeight))
		manager.wrandIndex[newIndex] = origin
		manager.quantityControl[origin] = 0
		manager.originMeta[origin] = meta
	}
}
func RemoveOrigin(label string, origin string) {

	if manager, ok := __tools[label]; ok {
		for index, ori := range manager.wrandIndex {
			if ori == origin {
				manager.wrand.SetWeight(index, uint(0))
				return
			}
		}
	}
}

func SetQuantityReporter(label string, reporter QuantityReporter) {

	if manager, ok := __tools[label]; ok {

		manager.reporter = reporter
	}
}
