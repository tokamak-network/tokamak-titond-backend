package tasks

import (
	"testing"
)

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
