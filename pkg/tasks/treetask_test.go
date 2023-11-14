package tasks

import (
	"errors"
	"testing"
)

func TestSubTasks(t *testing.T) {
	root := NewBaseTask(func(a int, b int, c int) (int, error) {
		return a + b + c, nil
	}, nil)

	task1 := NewBaseTask(func(a int) (int, error) {
		return a << 1, nil
	}, nil)

	task2 := NewBaseTask(func(a int) (int, error) {
		return a + 1000, nil
	}, nil)
	task3 := NewBaseTask(func(a int) (int, error) {
		return a + 12, nil
	}, nil)
	task4 := NewBaseTask(func(a int) (int, error) {
		return a + 24, nil
	}, nil)

	treeTask := NewTreeTask(root)
	treeTask.AddTask(task1, root)
	treeTask.AddTask(task2, root)
	treeTask.AddTask(task3, task2)
	treeTask.AddTask(task4, task2)

	if len(treeTask.Mapper[root]) != 2 {
		t.Error("Add subtask failed at root")
	}

	rootSubTasks := treeTask.Mapper[root]
	if rootSubTasks[0] != task1 {
		t.Error("Root: subtask1 failed")
	}
	if rootSubTasks[1] != task2 {
		t.Error("Root: subtask2 failed")
	}

	if len(treeTask.Mapper[task1]) != 0 {
		t.Error("Add subtask failed at task 1")
	}

	if len(treeTask.Mapper[task2]) != 2 {
		t.Error("Add subtask failed at task 2")
	}

	task2SubTasks := treeTask.Mapper[task2]
	if task2SubTasks[0] != task3 {
		t.Error("Task 2: subtask1 failed")
	}
	if task2SubTasks[1] != task4 {
		t.Error("Task 2: subtask2 failed")
	}

	if len(treeTask.Mapper[task3]) != 0 {
		t.Error("Add subtask failed at task 2")
	}

	if len(treeTask.Mapper[task4]) != 0 {
		t.Error("Add subtask failed at task 2")
	}
}

func TestNormalCases(t *testing.T) {
	root := NewBaseTask(func(a int, b int, c int) (int, error) {
		return a + b + c, nil
	}, func(result int, e error) {
		if result != 8 {
			t.Error(" TestFunctions: Root should return 8")
		}
		if e != nil {
			t.Error(" TestFunctions: Root should not return error")
		}
	})

	task1 := NewBaseTask(func(a int) (int, error) {
		return a << 1, nil
	}, func(result int, e error) {
		if result != 16 {
			t.Error(" TestFunctions: Task1 should return 16")
		}
		if e != nil {
			t.Error(" TestFunctions: Task1 should not return error")
		}
	})

	task2 := NewBaseTask(func(a int) (int, error) {
		return a + 1000, nil
	}, func(result int, e error) {
		if result != 1008 {
			t.Error(" TestFunctions: Task2 should return 1008")
		}
		if e != nil {
			t.Error(" TestFunctions: Task2 should not return error")
		}
	})

	task3 := NewBaseTask(func(a int) (int, error) {
		return a - 16, nil
	}, func(result int, e error) {
		if result != 0 {
			t.Error(" TestFunctions: Task3 should return 1008")
		}
		if e != nil {
			t.Error(" TestFunctions: Task3 should not return error")
		}
	})

	task4 := NewBaseTask(func(a int) (int, error) {
		return a + 16, nil
	}, func(result int, e error) {
		if result != 32 {
			t.Error(" TestFunctions: Task4 should return 1008")
		}
		if e != nil {
			t.Error(" TestFunctions: Task4 should not return error")
		}
	})

	task5 := NewBaseTask(func(a int) (int, error) {
		return a - 8, nil
	}, func(result int, e error) {
		if result != 1000 {
			t.Error(" TestFunctions: Task5 should return 1000")
		}
		if e != nil {
			t.Error(" TestFunctions: Task5 should not return error")
		}
	})

	task6 := NewBaseTask(func(a int) (int, int, error) {
		return a - 8, a, nil
	}, func(updatedNumber int, originalNumber int, e error) {
		if updatedNumber != 1000 {
			t.Error(" TestFunctions: Task5: updatedNumber should be 1000")
		}
		if originalNumber != 1008 {
			t.Error(" TestFunctions: Task5: originalNumber should be 1008")
		}
		if e != nil {
			t.Error(" TestFunctions: Task5 should not return error")
		}
	})

	treeTask := NewTreeTask(root)
	treeTask.AddTask(task1, root)
	treeTask.AddTask(task2, root)
	treeTask.AddTask(task3, task1)
	treeTask.AddTask(task4, task1)
	treeTask.AddTask(task5, task2)
	treeTask.AddTask(task6, task2)

	treeTask.Run(1, 2, 5)
}

func TestAbnormalCases(t *testing.T) {
	root := NewBaseTask(func(a int, b int, c int) (int, error) {
		return a + b + c, nil
	}, func(result int, e error) {
		if result != 8 {
			t.Error(" TestFunctions: Root should return 8")
		}
		if e != nil {
			t.Error(" TestFunctions: Root should not return error")
		}
	})

	task1 := NewBaseTask(func(a int) (int, error) {
		return a << 1, nil
	}, func(result int, e error) {
		if result != 16 {
			t.Error(" TestFunctions: Task1 should return 16")
		}
		if e != nil {
			t.Error(" TestFunctions: Task1 should not return error")
		}
	})

	task2 := NewBaseTask(func(a int) (int, error) {
		return 0, errors.New("")
	}, func(result int, e error) {
		if e == nil {
			t.Error(" TestFunctions: Task2 should return error")
		}
	})

	task3 := NewBaseTask(func(a int) (int, error) {
		return a - 16, nil
	}, func(result int, e error) {
		if result != 0 {
			t.Error(" TestFunctions: Task3 should return 1008")
		}
		if e != nil {
			t.Error(" TestFunctions: Task3 should not return error")
		}
	})

	task4 := NewBaseTask(func(a int) (int, error) {
		return a + 16, nil
	}, func(result int, e error) {
		if result != 32 {
			t.Error(" TestFunctions: Task4 should return 1008")
		}
		if e != nil {
			t.Error(" TestFunctions: Task4 should not return error")
		}
	})

	task5 := NewBaseTask(func(a int) (int, error) {
		return a - 8, nil
	}, func(result int, e error) {
		t.Error(" TestFunctions: Task5 callback should not be called")
	})

	treeTask := NewTreeTask(root)
	treeTask.AddTask(task1, root)
	treeTask.AddTask(task2, root)
	treeTask.AddTask(task3, task1)
	treeTask.AddTask(task4, task1)
	treeTask.AddTask(task5, task2)

	treeTask.Run(1, 2, 5)
}
