package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrInvalidCountWorkers = errors.New("invalid count of workers")

type Task func() error

// Запуск воркеров.
func startWorkers(wg *sync.WaitGroup, chTask <-chan Task, chErr chan<- error, done chan<- bool, stop <-chan bool) {
	defer wg.Done()

Loop:
	for {
		select {
		case <-stop:
			break Loop
		case task := <-chTask:
			err := task()
			chErr <- err
		}
	}
	done <- true
}

func startProduser(tasks []Task, wg *sync.WaitGroup, chTask chan<- Task, stop <-chan bool) {
	defer wg.Done()
	for _, task := range tasks {
		select {
		case <-stop:
			return
		case chTask <- task:
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

	var chTask = make(chan Task) // канал для выдачи задач.
	var chErr = make(chan error) // канал для приема результатов выполнения задач.
	var stop = make(chan bool)
	var done = make(chan bool, n)
	defer close(chTask)
	defer close(chErr)
	defer close(done)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go startWorkers(&wg, chTask, chErr, done, stop)
	}

	wg.Add(1)
	go startProduser(tasks, &wg, chTask, stop)
	// В цикле читаем результаты выполнения задач из канала chErr,
	// подсчитываем кол-во ошибок и выполненных задач.
	countErr := 0
	doneWorkers := 0
	doneTasks := 0
	flag := false
Loop:
	for {
		select {
		case err := <-chErr:
			if err != nil {
				countErr++
			} else {
				doneTasks++
			}
			if (countErr >= m || countErr >= len(tasks) || doneTasks >= len(tasks)) && !flag {
				close(stop)
				flag = true
			}
		case <-done:
			doneWorkers++
			if doneWorkers >= n {
				break Loop
			}
		}
	}

	wg.Wait()

	// Возврат ошибки, если кол-во задач с ошибками >= m.
	if countErr >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
