package hw06_pipeline_execution //nolint:golint,stylecheck
import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	wg := sync.WaitGroup{}

	stageFn := func(chIn In) {
		defer wg.Done()

		resultStage := chIn
		// запуск переданных stages.
		for _, stage := range stages {
			resultStage = stage(resultStage)
		}

		select {
		case <-done:
		case result := <-resultStage: // ждем результат выполнения stages.
			select {
			case <-done:
			case out <- result:
			}
		}
	}

	// читаем значения из канала in и запускаем обработчики
	go func() {
		if len(stages) != 0 {
			for num := range in {
				chIn := make(Bi)
				wg.Add(1)
				go stageFn(chIn)
				select {
				case <-done:
				case chIn <- num:
				}
			}
		}
		wg.Wait()
		close(out)
	}()

	return out
}
