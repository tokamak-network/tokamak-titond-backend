package tasks

import (
	"errors"
	"reflect"
)

type ITaskResult struct {
	Response interface{}
	Err      interface{}
}

type ITask interface {
	Exec() (interface{}, error)
	AddHandler(interface{})
	AddInput([]interface{})
}

type BaseTask struct {
	handler  interface{}
	input    []interface{}
	callback interface{}
}

func (task *BaseTask) AddHandler(handler interface{}) {
	task.handler = handler
}

func (task *BaseTask) AddInput(input []interface{}) {
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
	reflect.ValueOf(task.callback).Call(results)
	return results, nil
}

func buildCallParams(a ...interface{}) []reflect.Value {
	params := make([]reflect.Value, len(a))
	for i, value := range a {
		params[i] = reflect.ValueOf(value)
	}
	return params
}
