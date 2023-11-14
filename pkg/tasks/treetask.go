package tasks

import (
	"fmt"
	"reflect"
)

type TreeTask struct {
	Root     ITask
	Mapper   map[ITask][]ITask
	InputMap map[ITask][]reflect.Value
}

func (task *TreeTask) AddTask(child ITask, parent ITask) {
	subTasks, exist := task.Mapper[parent]
	if !exist {
		task.Mapper[parent] = make([]ITask, 0)
	}
	task.Mapper[parent] = append(subTasks, child)
}

func (task *TreeTask) Run(input ...interface{}) {
	processingTasks := make([]ITask, 0)
	processingTasks = append(processingTasks, task.Root)
	params := buildCallParams(input...)
	task.updateInputMap(task.Root, params)
	for len(processingTasks) > 0 {
		currentTask := processingTasks[0]
		params = task.getInputMap(currentTask)
		params, err := currentTask.LowLevelCall(params)
		if err == nil {
			nextTasks := task.Mapper[currentTask]
			for _, iTask := range nextTasks {
				task.updateInputMap(iTask, params)
			}
			processingTasks = append(processingTasks, nextTasks...)
		}

		processingTasks = processingTasks[1:]
	}
}

func (task *TreeTask) Interate() {
	processingTasks := make([]ITask, 0)
	processingTasks = append(processingTasks, task.Root)
	for len(processingTasks) > 0 {
		currentTask := processingTasks[0]
		fmt.Println(currentTask)
		nextTasks := task.Mapper[currentTask]
		processingTasks = append(processingTasks, nextTasks...)
		processingTasks = processingTasks[1:]
	}
}

func (task *TreeTask) updateInputMap(iTask ITask, input []reflect.Value) {
	task.InputMap[iTask] = input
}

func (task *TreeTask) getInputMap(iTask ITask) []reflect.Value {
	if input, exist := task.InputMap[iTask]; exist {
		return input
	}
	return []reflect.Value{}
}

func NewTreeTask(root ITask) *TreeTask {
	mapper := make(map[ITask][]ITask)
	mapper[root] = make([]ITask, 0)
	inputMap := make(map[ITask][]reflect.Value)
	return &TreeTask{Root: root, Mapper: mapper, InputMap: inputMap}
}
