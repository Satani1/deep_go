package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	tasks   []Task
	taskMap map[int]int
}

func NewScheduler() Scheduler {
	return Scheduler{
		tasks:   make([]Task, 0),
		taskMap: make(map[int]int, 0),
	}
}

func (s *Scheduler) AddTask(task Task) {
	if _, exists := s.taskMap[task.Identifier]; exists {
		return
	}

	s.tasks = append(s.tasks, task)
	index := len(s.tasks) - 1
	s.taskMap[task.Identifier] = index
	s.up(index)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	index, exists := s.taskMap[taskID]
	if !exists {
		return
	}

	oldPriority := s.tasks[index].Priority
	s.tasks[index].Priority = newPriority

	if newPriority > oldPriority {
		s.up(index)
	} else if newPriority < oldPriority {
		s.down(index)
	}
}

func (s *Scheduler) GetTask() Task {
	if len(s.tasks) == 0 {
		return Task{}
	}

	root := s.tasks[0]
	last := len(s.tasks) - 1
	s.tasks[0] = s.tasks[last]
	s.taskMap[s.tasks[0].Identifier] = 0
	s.tasks = s.tasks[:last]
	delete(s.taskMap, root.Identifier)

	if len(s.tasks) > 0 {
		s.down(0)
	}

	return root
}

func (s *Scheduler) up(i int) {
	for {
		parent := (i - 1) / 2
		if parent == i || s.tasks[parent].Priority >= s.tasks[i].Priority {
			break
		}
		s.swap(i, parent)
		i = parent
	}
}

func (s *Scheduler) down(i int) {
	for {
		left := 2*i + 1
		if left >= len(s.tasks) || left < 0 {
			break
		}
		largest := left
		if right := left + 1; right < len(s.tasks) &&
			s.tasks[right].Priority > s.tasks[left].Priority {
			largest = right
		}
		if s.tasks[i].Priority >= s.tasks[largest].Priority {
			break
		}
		s.swap(i, largest)
		i = largest
	}
}

func (s *Scheduler) swap(i, j int) {
	s.tasks[i], s.tasks[j] = s.tasks[j], s.tasks[i]
	s.taskMap[s.tasks[i].Identifier] = i
	s.taskMap[s.tasks[j].Identifier] = j
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 1, Priority: 100}, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
