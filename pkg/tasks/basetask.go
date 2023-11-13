package tasks

import (
	"errors"
	"fmt"
	"reflect"
)

type BaseTask struct {
	handler  interface{}
	input    []interface{}
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

func (task *BaseTask) Exec() (interface{}, error) {
	params := buildCallParams(task.input...)
	results := reflect.ValueOf(task.handler).Call(params)
	if len(results) != 2 {
		return nil, errors.New("the task should have a tuple of 2 items as returned item")
	}
	if task.callback != nil {
		reflect.ValueOf(task.callback).Call(results)
	}
	return results, nil
}

func (task *BaseTask) Run(input ...interface{}) (interface{}, error) {

	params := buildCallParams(input...)
	results := reflect.ValueOf(task.handler).Call(params)
	if len(results) != 2 {
		return nil, errors.New("the task should have a tuple of 2 items as returned item")
	}
	if task.callback != nil {
		reflect.ValueOf(task.callback).Call(results)
	}
	return results, nil
}

func MakeBaseTask(handler interface{}, callback interface{}) *BaseTask {
	return &BaseTask{handler: handler, callback: callback}
}

func buildCallParams(a ...interface{}) []reflect.Value {
	fmt.Println("Check length: ", len(a))
	params := make([]reflect.Value, len(a))
	for i, value := range a {
		params[i] = reflect.ValueOf(value)
	}
	return params
}
