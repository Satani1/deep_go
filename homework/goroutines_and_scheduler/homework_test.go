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

type TaskHeap struct {
	tasks       []Task
	taskIndexes map[int]int
}

func NewTaskHeap() *TaskHeap {
	return &TaskHeap{
		tasks:       []Task{},
		taskIndexes: make(map[int]int),
	}
}

func (t *TaskHeap) Len() int {
	return len(t.tasks)
}

func (t *TaskHeap) up(i int) {
	for {
		parent := (i - 1) / 2
		if parent == i || t.tasks[parent].Priority >= t.tasks[i].Priority {
			break
		}

		t.swap(i, parent)
		i = parent
	}
}

func (t *TaskHeap) down(i int) {
	for {
		left := 2*i + 1
		if left >= t.Len() || left < 0 {
			break
		}

		largest := left
		if right := left + 1; right < t.Len() && t.tasks[right].Priority > t.tasks[left].Priority {
			largest = right
		}

		if t.tasks[i].Priority >= t.tasks[largest].Priority {
			break
		}

		t.swap(i, largest)
		i = largest
	}
}

func (t *TaskHeap) swap(i, j int) {
	t.tasks[i], t.tasks[j] = t.tasks[j], t.tasks[i]
	t.taskIndexes[t.tasks[i].Identifier] = i
	t.taskIndexes[t.tasks[j].Identifier] = j
}

func (t *TaskHeap) Push(task Task) {
	if _, exists := t.taskIndexes[task.Identifier]; exists {
		return
	}

	t.tasks = append(t.tasks, task)
	index := t.Len() - 1
	t.taskIndexes[task.Identifier] = index
	t.up(index)
}

func (t *TaskHeap) Pop() Task {
	if t.Len() == 0 {
		return Task{}
	}

	root := t.tasks[0]
	last := t.Len() - 1
	t.tasks[0] = t.tasks[last]
	t.taskIndexes[t.tasks[0].Identifier] = 0
	t.tasks = t.tasks[:last]
	delete(t.taskIndexes, root.Identifier)

	if t.Len() > 0 {
		t.down(0)
	}

	return root
}

func (t *TaskHeap) Update(identifier int, priority int) {
	index, exists := t.taskIndexes[identifier]
	if !exists {
		return
	}

	oldPriority := t.tasks[index].Priority
	t.tasks[index].Priority = priority

	if priority > oldPriority {
		t.up(index)
	} else if priority < oldPriority {
		t.down(index)
	}
}

type Scheduler struct {
	taskHeap *TaskHeap
}

func NewScheduler() Scheduler {
	return Scheduler{
		taskHeap: NewTaskHeap(),
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.taskHeap.Push(task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.taskHeap.Update(taskID, newPriority)
}

func (s *Scheduler) GetTask() Task {
	return s.taskHeap.Pop()
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
