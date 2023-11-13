package tasks

import (
	"testing"
)

func TestRunTask(t *testing.T) {
	adder := func(a int, b int) (int, error) {
		return (a + b), nil
	}

	task := MakeBaseTask(adder, func(result int, err error) {
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

	task := MakeBaseTask(adder, func(result int, err error) {
		if result != 6 {
			t.Error("Need to be equal")
		}
	})
	task.AddInput(3, 3)

	task.Exec()
}
