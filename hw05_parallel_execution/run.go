package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrInvalidCountWorkers = errors.New("invalid count of workers")

type Task func() error

// Запуск воркеров.
func startWorkers(wg *sync.WaitGroup, chTask <-chan Task, chErr chan<- error) {
	defer wg.Done()

	for task := range chTask {
		err := task()
		chErr <- err
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
//nolint:gocognit
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	if n <= 0 {
		return ErrInvalidCountWorkers
	}

	var chTask = make(chan Task, 1) // Канал для выдачи задач.
	var chErr = make(chan error)    // Канал для приема результатов выполнения задач.
	defer close(chErr)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go startWorkers(&wg, chTask, chErr)
	}

	countTasksOK := 0  // Выполнено задач успешно.
	countTasksErr := 0 // Выполнено задач с ошибкой.
	issuedTasks := 0   // Выдано задач.
	flag := false
	for (countTasksOK + countTasksErr) < len(tasks) {
		select {
		case result := <-chErr: // прием результатов работы воркера.
			if result == nil {
				countTasksOK++
			} else {
				countTasksErr++
			}
			if (countTasksErr >= m || issuedTasks == len(tasks)) && !flag {
				flag = true
				close(chTask)
			}
		default: // Выдача задач.
			if len(chTask) == 0 && issuedTasks < len(tasks) && countTasksErr < m {
				chTask <- tasks[issuedTasks]
				issuedTasks++
			}
		}
		// все выданные задачи выполнены?
		if countTasksOK+countTasksErr == issuedTasks && countTasksErr > 1 {
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
