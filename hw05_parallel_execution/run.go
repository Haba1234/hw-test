package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrInvalidCountWorkers = errors.New("invalid count of workers")

type Task func() error

// Запуск воркеров.
func startWorkers(wg *sync.WaitGroup, chTask <-chan Task, chErr chan<- error, done <-chan struct{}) {
	defer wg.Done()

	for task := range chTask {
		select {
		case <-done:
			return
		case chErr <- task():
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	if n <= 0 {
		return ErrInvalidCountWorkers
	}

	var chTask = make(chan Task) // Канал для выдачи задач.
	var chErr = make(chan error) // Канал для приема результатов выполнения задач.
	var done = make(chan struct{})
	defer close(chErr)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go startWorkers(&wg, chTask, chErr, done)
	}

	countTasksErr := 0 // Выполнено задач с ошибкой.
	issuedTasks := 0   // Выдано задач.

	for {
		select {
		case result := <-chErr: // прием результатов работы воркера.
			if result != nil {
				countTasksErr++
			}
		case chTask <- tasks[issuedTasks]: // Раздача новых задач.
			issuedTasks++
		}
		if countTasksErr >= m || issuedTasks == len(tasks) {
			close(chTask)
			close(done)
			break
		}
	}

	wg.Wait()

	// Возврат ошибки, если кол-во задач с ошибками >= m.
	if countTasksErr >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
