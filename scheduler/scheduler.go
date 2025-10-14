package scheduler

import "sync"

// Task represents a function to be executed.
type Task func()

// TaskRunner manages the concurrent execution of a set of tasks.
type TaskRunner struct {
	tasks []Task
}

// NewTaskRunner creates a new runner for concurrent tasks.
func NewTaskRunner(tasks ...Task) *TaskRunner {
	return &TaskRunner{
		tasks: tasks,
	}
}

// Run executes all registered tasks in parallel and waits for them to complete.
func (tr *TaskRunner) Run() {
	var wg sync.WaitGroup
	for _, task := range tr.tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			t()
		}(task)
	}
	wg.Wait()
}