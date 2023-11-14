package tasks

import (
	"testing"
)

func TestNoReturnItem(t *testing.T) {
	adder := func(a int) {
	}

	task := NewBaseTask(adder, func() {
	})
	_, err := task.Run(8)
	if err == nil {
		t.Error("Should raise error since does not have return value")
	}
}

func TestSingleParamTask(t *testing.T) {
	adder := func(a int) (int, error) {
		return (a + 1), nil
	}

	task := NewBaseTask(adder, func(result int, err error) {
		if result != 9 {
			t.Error("Need to be equal")
		}
	})
	task.Run(8)
}

func TestRunTask(t *testing.T) {
	adder := func(a int, b int) (int, error) {
		return (a + b), nil
	}

	task := NewBaseTask(adder, func(result int, err error) {
		if result != 5 {
			t.Error("Need to be equal")
		}
	})

	task.Run(2, 3)
}

func TestExecTask(t *testing.T) {
	adder := func(a int, b int) (int, error) {
		return (a + b), nil
	}

	task := NewBaseTask(adder, func(result int, err error) {
		if result != 6 {
			t.Error("Need to be equal")
		}
	})

	task.AddInput(3, 3)
	task.Exec()
}

func TestExecTaskWithMoreThan2ReturnField(t *testing.T) {
	operator := func(a int, b int) (int, int, error) {
		return a + b, a - b, nil
	}

	task := NewBaseTask(operator, func(add int, sub int, err error) {
		if add != 6 {
			t.Error("Add failed")
		}
		if sub != 0 {
			t.Error("Sub failed")
		}
	})

	task.AddInput(3, 3)
	task.Exec()
}

type Object struct {
	val int
}

func (obj *Object) Increase() error {
	obj.val++
	return nil
}

func (obj *Object) Set(val int) error {
	obj.val = val
	return nil
}

func (obj *Object) Add(a int, b int, c int) error {
	obj.val = a + b + c
	return nil
}

func TestWithStruct(t *testing.T) {
	obj := &Object{val: 5}
	task := NewBaseTask(obj.Increase, func(e error) {
		if e != nil {
			t.Error("e should be nil")
		}

		if obj.val != 6 {
			t.Error("Val should be 6")
		}

	})
	_, err := task.Run()
	if err != nil {
		t.Error("Should run task successfully")
	}

	task1 := NewBaseTask(obj.Set, func(e error) {
		if e != nil {
			t.Error("e should be nil")
		}

		if obj.val != 15 {
			t.Error("Val should be 15 after directly set")
		}

	})
	_, err = task1.Run(15)
	if err != nil {
		t.Error("Should run task 1 successfully")
	}

	task2 := NewBaseTask(obj.Add, func(e error) {
		if e != nil {
			t.Error("e should be nil")
		}

		if obj.val != 30 {
			t.Error("Val should be 30 after called add()")
		}

	})
	_, err = task2.Run(5, 10, 15)
	if err != nil {
		t.Error("Should run task 1 successfully")
	}
}
