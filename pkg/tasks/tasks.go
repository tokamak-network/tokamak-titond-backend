package tasks

import (
	"reflect"
)

type ITaskResult struct {
	Response interface{}
	Err      interface{}
}

type ITask interface {
	AddHandler(interface{})
	AddInput(...interface{})
	AddCallback(interface{})
	Handler() interface{}
	involkeCallback()
	Exec() (interface{}, error)
	Run(input ...interface{}) (interface{}, error)
	LowLevelCall([]reflect.Value) ([]reflect.Value, error)
}

func buildCallParams(a ...interface{}) []reflect.Value {
	params := make([]reflect.Value, len(a))
	for i, value := range a {
		params[i] = reflect.ValueOf(value)
	}
	return params
}
