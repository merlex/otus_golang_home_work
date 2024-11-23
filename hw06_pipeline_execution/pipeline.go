package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		c := make(Bi)
		go func() {
			defer close(c)
			if done != nil {
				<-done
				return
			}
		}()
		return c
	}

	resultChan := in
	for _, st := range stages {
		outStage := make(Bi)
		go func(stage Stage, chanIn In, chanOut Bi) {
			defer close(chanOut)
			currentChanOut := stage(chanIn)
			for {
				select {
				case d, ok := <-currentChanOut:
					if !ok {
						return
					}
					chanOut <- d
				case <-done:
					return
				}
			}
		}(st, resultChan, outStage)
		resultChan = outStage
	}
	return resultChan
}
