package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func runStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case result, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case out <- result:
					}
				}
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outIn := in
	for _, stage := range stages {
		outIn = stage(runStage(outIn, done))
	}
	return runStage(outIn, done)
}
