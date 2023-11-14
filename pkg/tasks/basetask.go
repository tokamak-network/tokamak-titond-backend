package tasks

import (
	"errors"
	"reflect"
)

type BaseTask struct {
	handler  interface{}
	input    []interface{}
	results  []reflect.Value
	callback interface{}
}

func (task *BaseTask) AddHandler(handler interface{}) {
	task.handler = handler
}

func (task *BaseTask) AddInput(input ...interface{}) {
	task.input = input
}

func (task *BaseTask) AddCallback(callback interface{}) {
	task.callback = callback
}

func (task *BaseTask) Handler() interface{} {
	return task.handler
}

func (task *BaseTask) involkeCallback() {
	if task.callback != nil {
		reflect.ValueOf(task.callback).Call(task.results)
	}
}

func (task *BaseTask) Exec() (interface{}, error) {
	params := buildCallParams(task.input...)
	task.results = reflect.ValueOf(task.handler).Call(params)
	if len(task.results) != 2 {
		return nil, errors.New("the task should have a tuple of 2 items as returned item")
	}
	task.involkeCallback()
	return task.results, nil
}

func (task *BaseTask) Run(input ...interface{}) (interface{}, error) {
	params := buildCallParams(input...)
	task.results = reflect.ValueOf(task.handler).Call(params)
	if len(task.results) != 2 {
		return nil, errors.New("the task should have a tuple of 2 items as returned item")
	}
	task.involkeCallback()
	return task.results, nil
}

func (task *BaseTask) LowLevelCall(params []reflect.Value) ([]reflect.Value, error) {
	results := reflect.ValueOf(task.handler).Call(params)
	task.results = results
	if err, ok := results[len(results)-1].Interface().(error); ok {
		return results[0 : len(results)-1], err
	} else {
		task.results = results
		task.involkeCallback()
		return results[0 : len(results)-1], nil
	}
}

func NewBaseTask(handler interface{}, callback interface{}) *BaseTask {
	return &BaseTask{handler: handler, callback: callback}
}
