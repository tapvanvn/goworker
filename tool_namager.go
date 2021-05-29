package goworker

import "math"

type ToolManager struct {
	BlackSmith IToolMaker
	numMade    int
	stack      chan interface{}
	max        int
}

var __tools = map[string]*ToolManager{}

func AddTool(label string, toolMaker IToolMaker) {

	AddToolWithControl(label, toolMaker, math.MaxInt32)
}
func AddToolWithControl(label string, toolMaker IToolMaker, maxTool int) {

	if manager, ok := __tools[label]; !ok {

		m := &ToolManager{
			BlackSmith: toolMaker,
			stack:      make(chan interface{}, 1),
			numMade:    1,
			max:        maxTool,
		}
		m.stack <- m.BlackSmith.Make()
		__tools[label] = m

	} else {

		manager.BlackSmith = toolMaker
	}
}
func borrow(label string) interface{} {

	if manager, ok := __tools[label]; ok {
		if len(manager.stack) == 0 && manager.numMade < manager.max {
			manager.numMade++
			return manager.BlackSmith.Make()
		}

		return <-manager.stack
	}
	return nil
}

func thankyou(label string, tool interface{}) {
	if manager, ok := __tools[label]; ok {

		manager.stack <- tool
	}
}
