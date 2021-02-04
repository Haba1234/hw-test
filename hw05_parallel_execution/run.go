package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrInvalidCountWorkers = errors.New("invalid count of workers")

type Task func() error

// Функция запуска воркеров.
func startGo(tasks []Task, wg *sync.WaitGroup, chNumTask chan int, chErr chan error) {
	defer wg.Done()
	for {
		num, ok := <-chNumTask // ждем номер следующей задачи.
		if !ok {
			return // Канал с задачами закрыт. Завершаем работу.
		}
		err := tasks[num]()
		chErr <- err // пишем результат в канал.
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	// Проверка, если m меньше 1, то возвращаем ошибку.
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	// Проверка, если n меньше 1, то возвращаем ошибку.
	if n <= 0 {
		return ErrInvalidCountWorkers
	}
	// Открываем канал для приема результатов выполнения задач.
	var chErr = make(chan error, n)
	defer close(chErr)

	// Открываем канал для выдачи номеров задач.
	var chNumTask = make(chan int, n)

	currentTask := 0 // Подсчет кол-ва выполненных задач.

	wg := sync.WaitGroup{}
	wg.Add(n)
	// Запускаем n воркеров.
	for i := 0; i < n; i++ {
		go startGo(tasks, &wg, chNumTask, chErr)
	}

	countErr := 0 // Кол-во задач с ошибками.
	// В цикле читаем результаты выполнения задач из канала chErr,
	// подсчитываем кол-во ошибок и
	// раздаем новые задачи освободившимся воркерам через канал currentTask.
	for currentTask < len(tasks) {
		chNumTask <- currentTask // Новая задача.
		currentTask++
		if err := <-chErr; true {
			if err != nil {
				countErr++
			}
			if countErr >= m { // Достигнут предел максимального кол-ва ошибок.
				break
			}
		}
	}
	// Подаем сигнал воркерам на прекращение работы.
	close(chNumTask)

	wg.Wait() // Ожидаем завершение всех запущенных воркеров.

	// Возврат ошибки, если кол-во задач с ошибками >= m.
	if countErr >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
