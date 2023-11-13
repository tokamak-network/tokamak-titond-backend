package tasks

type ITaskResult struct {
	Response interface{}
	Err      interface{}
}

type ITask interface {
	Exec() (interface{}, error)
	AddHandler(interface{})
	AddInput([]interface{})
}
